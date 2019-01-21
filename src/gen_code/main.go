package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type method struct {
	index         int
	name          string
	req_type_name string
	rsp_type_name string
}

type service struct {
	http bool
	name string
	ms   []*method
}

type pb struct {
	name string
	ss   []*service
}

func error_check(err error) {
	if err != nil {
		panic(err)
	}
}

func handler_http(p *pb, orgi_file_name string) {
	for _, s := range p.ss {
		if !s.http {
			continue
		}
		suffix := fmt.Sprintf(".%s.http.go", s.name)
		fn := strings.Replace(orgi_file_name, ".proto", suffix, 1)
		fmt.Println("gen http filename: ", fn)

		f, err := os.OpenFile(fn, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		error_check(err)
		defer f.Close()
		defer f.Sync()

		// package
		{
			f.WriteString(fmt.Sprintf("package %s\n", p.name))
			f.WriteString("\n")
		}

		svr_name := s.name // strings.ToLower(s.name)
		// pb_str := fmt.Sprintf("%s_pb", svr_name)

		// import
		{
			f.WriteString("import (\n")
			f.WriteString("\t\"context\"\n")
			f.WriteString("\t. \"reflect\"\n")
			f.WriteString("\t\"encoding/json\"\n")
			f.WriteString("\t. \"github.com/loveCupid/dipamkara/src/kernal\"\n")
			f.WriteString(")\n")
			f.WriteString("\n")
		}

		// struct
		{
			f.WriteString(fmt.Sprintf("type %s_http struct {\n", svr_name))
			f.WriteString(fmt.Sprintf("\tc %sClient\n", svr_name))
			f.WriteString("}\n")
			f.WriteString("\n")
		}

		// run method
		{
			run_method_str :=
				`
func Run%sHttp() {
	s  := NewServer("%s_http")
	sh := new(%s_http)

    sc, err := FetchServiceConn("%s", s)
    ErrorPanic(err)
	sh.c = New%sClient(sc)

	RegisterHttpServiceServer(s.Svr, sh)
	s.Svr.Serve(s.Lis)
}
`
			f.WriteString(fmt.Sprintf(run_method_str,
				svr_name, svr_name, svr_name,
				svr_name, svr_name,
			))
		}

		// Call method
		{
			call_method_str := `
func (s *%s_http) Call(ctx context.Context, in *HttpReq) (*HttpRsp, error) {

    val := TypeOf(s)
    m,ok := val.MethodByName("Call" + in.Method)
    if !ok {
        Error(ctx, "not found method. method: ", in.Method)
        return nil, nil
    }

    resp := m.Func.Call([]Value{ValueOf(s), ValueOf(ctx), ValueOf(in)})

    // return resp[0].Interface().(*HttpRsp), resp[1].Interface().(error)
    return resp[0].Interface().(*HttpRsp), nil 
}
            `
			f.WriteString(fmt.Sprintf(call_method_str, svr_name))
			f.WriteString("\n")
		}

		// method
		for _, m := range s.ms {
			call_method_str := `
func (s *%s_http) Call%s(ctx context.Context, in *HttpReq) (*HttpRsp, error) {

    var req %s
    json.Unmarshal([]byte(in.Body), &req)

    resp,err := s.c.%s(ctx, &req)
    resp_json_str, err := json.Marshal(resp)
    if err != nil {
        return nil, err
    }

    return &HttpRsp{Reply: string(resp_json_str)}, nil
}
            `
			f.WriteString(fmt.Sprintf(call_method_str, svr_name, m.name, m.req_type_name, m.name))
			f.WriteString("\n")
		}
	}
}

func handler_svr(p *pb, orgi_file_name string) {
	for _, s := range p.ss {
		suffix := fmt.Sprintf(".%s.svr.go", s.name)
		fn := strings.Replace(orgi_file_name, ".proto", suffix, 1)

		farr := strings.Split(fn, "/")
		farr[len(farr)-1] = "_" + farr[len(farr)-1]
		fn = strings.Join(farr, "/")

		fmt.Println("gen svr filename: ", fn)

		f, err := os.OpenFile(fn, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		error_check(err)
		defer f.Close()
		defer f.Sync()

		// package
		{
			f.WriteString("package main\n")
			f.WriteString("\n")
		}

		svr_name := s.name // strings.ToLower(s.name)
		pb_str := fmt.Sprintf("%s_pb", svr_name)

		// import
		{
			f.WriteString("import (\n")
			f.WriteString("\t\"context\"\n")
			f.WriteString("\t. \"github.com/loveCupid/dipamkara/src/kernal\"\n")
			f.WriteString("\t" + pb_str + fmt.Sprintf(" \"github.com/loveCupid/dipamkara/src/%s/proto\"\n", p.name))
			f.WriteString(")\n")
			f.WriteString("\n")
		}

		// struct
		{
			f.WriteString(fmt.Sprintf("type %s_svr struct {}\n", svr_name))
			f.WriteString("\n")
		}

		// run method
		{
			f.WriteString(fmt.Sprintf("func Run%s() {\n", svr_name))
			f.WriteString(fmt.Sprintf("\ts := NewServer(\"%s\")\n", svr_name))
			f.WriteString(fmt.Sprintf("\t%s.Register%sServer(s.Svr, &%s_svr{})\n", pb_str, svr_name, svr_name))
			if s.http {
				f.WriteString(fmt.Sprintf("\tgo %s.Run%sHttp()\n", pb_str, svr_name))
			}
			f.WriteString("\ts.Svr.Serve(s.Lis)\n")
			f.WriteString("}\n")
			f.WriteString("\n")
		}

		// method
		for _, m := range s.ms {
			f.WriteString(fmt.Sprintf("func (s *%s_svr) %s(ctx context.Context, in *%s.%s) (*%s.%s, error) {\n",
				svr_name, m.name, pb_str, m.req_type_name, pb_str, m.rsp_type_name))
			f.WriteString("\treturn nil, nil\n")
			f.WriteString("}\n")
			f.WriteString("\n")
		}
	}
}
func handler_cli(p *pb, orgi_file_name string) {
	for _, s := range p.ss {
		suffix := fmt.Sprintf(".%s.cli.go", s.name)
		fn := strings.Replace(orgi_file_name, ".proto", suffix, 1)
		fmt.Println("gen cli filename: ", fn)

		f, err := os.OpenFile(fn, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		error_check(err)
		defer f.Close()
		defer f.Sync()

		// package
		{
			f.WriteString(fmt.Sprintf("package %s\n", p.name))
			f.WriteString("\n")
		}

		svr_name := s.name // strings.ToLower(s.name)
		// pb_str := fmt.Sprintf("%s_pb", svr_name)

		// import
		{
			f.WriteString("import (\n")
			f.WriteString("\t\"context\"\n")
			f.WriteString("\t. \"github.com/loveCupid/dipamkara/src/kernal\"\n")
			f.WriteString(")\n")
			f.WriteString("\n")
		}

		// method
		for _, m := range s.ms {
			method_str :=
				`
func Call_%s_%s(ctx context.Context, in *%s) (*%s, error) {
    sc, err := FetchServiceConnByCtx(ctx, "%s")
    if err != nil {
        Error(ctx, "fetch %s service conn error. ")
        return nil, err
    }
    return New%sClient(sc).%s(ctx, in)
}
            `
			f.WriteString(fmt.Sprintf(method_str,
				svr_name, m.name, m.req_type_name, m.rsp_type_name,
				svr_name, svr_name, svr_name, m.name,
			))
		}
	}
}

func handler(p *pb) {

	orgi_file_name := os.Args[1]

	fmt.Println("package name: ", p.name)
	for _, s := range p.ss {
		fmt.Println("\tservice: ", s.name, "http: ", s.http)
		for _, m := range s.ms {
			fmt.Println("\t\t", m.name, "(", m.req_type_name, ",", m.rsp_type_name, ")")
		}
	}

	// gen svr
	{
		handler_svr(p, orgi_file_name)
	}
	// gen cli
	{
		handler_cli(p, orgi_file_name)
	}
	// gen http
	{
		handler_http(p, orgi_file_name)
	}
}

func main() {
	if len(os.Args) == 1 {
		fmt.Printf("not found proto file...\n")
		os.Exit(0)
	}

	filename := os.Args[1]

	f, err := os.Open(filename)
	error_check(err)
	defer f.Close()

	p := new(pb)

	has_http := false

	rd := bufio.NewReader(f)
	for {
		line, err := rd.ReadString('\n')
		if err == io.EOF {
			fmt.Printf("finish!\n")
			break
		}
		error_check(err)

		line = strings.Trim(line, " ")

		// package name
		if strings.HasPrefix(line, "package ") {
			s := strings.Trim(line, "package")
			s = strings.Trim(s, "\n")
			s = strings.Trim(s, " ")
			s = strings.Trim(s, ";")
			s = strings.Trim(s, " ")
			p.name = s
		}

		// has http
		if strings.HasPrefix(line, "// outside") {
			has_http = true
		}

		// service
		if strings.HasPrefix(line, "service ") {
			_service := new(service)
			s := strings.Trim(line, "service")
			s = strings.Trim(s, "\n")
			s = strings.Trim(s, " ")
			s = strings.Trim(s, "{")
			s = strings.Trim(s, " ")
			_service.http = has_http
			_service.name = s

			p.ss = append(p.ss, _service)

			has_http = false
		}
		// method
		if strings.HasPrefix(line, "rpc ") {
			_service := p.ss[len(p.ss)-1]
			_method := new(method)
			_method.index = len(_service.ms)

			s := strings.Trim(line, "rpc")
			s = strings.Trim(s, " ")

			i := 0
			ii := strings.Index(s, "(")
			_method.name = string(s[i:ii])

			// i  = strings.Index(s[ii+1:], "(")
			i = ii + 1
			ii = i + strings.Index(s[i:], ")")
			fmt.Printf("i: %d, ii: %d\n", i, ii)
			_method.req_type_name = strings.Trim(s[i:ii], " ")

			i = strings.Index(s[ii:], "(") + ii + 1
			ii = strings.Index(s[i:], ")") + i
			fmt.Printf("i: %d, ii: %d\n", i, ii)
			_method.rsp_type_name = strings.Trim(s[i:ii], " ")
			_service.ms = append(_service.ms, _method)
		}
	}

	handler(p)
}
