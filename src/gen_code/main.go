package main
import (
    "io"
    "os"
    "fmt"
    "bufio"
    "strings"
)

type method struct {
    index int
    name string
    req_type_name string
    rsp_type_name string
}

type service struct {
    name string
    ms []*method
}

type pb struct {
    name string
    ss []*service
}

func error_check(err error) {
    if err != nil {
        panic(err)
    }
}

func handler_svr(p *pb, orgi_file_name string) {
    for _,s := range p.ss {
        suffix := fmt.Sprintf(".%s.svr.go", s.name)
        fn := strings.Replace(orgi_file_name, ".proto", suffix, 1)
        fmt.Println("gen svr filename: ", fn)

        f,err := os.OpenFile(fn, os.O_WRONLY | os.O_CREATE, 0666)
        error_check(err)
        defer f.Close()
        defer f.Sync()

        // package 
        {
            f.WriteString("package main\n")
            f.WriteString("\n")
        }

        svr_name := s.name// strings.ToLower(s.name)
        pb_str := fmt.Sprintf("%s_pb", svr_name)

        // import 
        {
            f.WriteString("import (\n")
            f.WriteString("\t\"context\"\n")
            f.WriteString("\t. \"github.com/loveCupid/dipamkara/src/kernal\"\n")
            f.WriteString("\t" + pb_str + " \"github.com/loveCupid/dipamkara/src/hello/proto\"\n")
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
            f.WriteString("\ts.Svr.Serve(s.Lis)\n")
            f.WriteString("}\n")
            f.WriteString("\n")
        }

        // method
        for _,m := range s.ms {
            f.WriteString(fmt.Sprintf("func (s *%s_svr) %s(ctx context.Context, in *%s.%s) (*%s.%s, error) {\n",
                svr_name, m.name, pb_str, m.req_type_name, pb_str, m.rsp_type_name))
            f.WriteString("\treturn nil, nil\n")
            f.WriteString("}\n")
            f.WriteString("\n")
        }
    }
}
func handler_cli(p *pb, orgi_file_name string) {
    for _,s := range p.ss {
        suffix := fmt.Sprintf(".%s.cli.go", s.name)
        fn := strings.Replace(orgi_file_name, ".proto", suffix, 1)
        fmt.Println("gen cli filename: ", fn)

        f,err := os.OpenFile(fn, os.O_WRONLY | os.O_CREATE, 0666)
        error_check(err)
        defer f.Close()
        defer f.Sync()

        // package 
        {
            f.WriteString(fmt.Sprintf("package %s\n", p.name))
            f.WriteString("\n")
        }

        svr_name := s.name// strings.ToLower(s.name)
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
        for _,m := range s.ms {
            f.WriteString(fmt.Sprintf("func Call_%s_%s(ctx context.Context, in *%s) (*%s, error) {\n",
                svr_name, m.name, m.req_type_name, m.rsp_type_name))
            f.WriteString(fmt.Sprintf("\tc := New%sClient(FetchServiceConnByCtx(ctx, \"%s\"))\n",
                svr_name, svr_name))
            f.WriteString(fmt.Sprintf("\treturn c.%s(ctx, in)\n",m.name))
            f.WriteString("}\n")
            f.WriteString("\n")
        }
    }
}

func handler(p *pb) {

    orgi_file_name := os.Args[1]

    fmt.Println("package name: ", p.name)
    for _,s := range p.ss {
        fmt.Println("\tservice: ", s.name)
        for _,m := range s.ms {
            fmt.Println("\t\t", m.name, "(", m.req_type_name, ",", m.rsp_type_name , ")")
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
}

func main() {
    if len(os.Args) == 1 {
        fmt.Printf("not found proto file...\n")
        os.Exit(0)
    }

    filename := os.Args[1]

    f,err := os.Open(filename)
    error_check(err)
    defer f.Close()

    p := new(pb)

    rd := bufio.NewReader(f)
    for {
        line,err := rd.ReadString('\n')
        if err == io.EOF {
            fmt.Printf("finish!\n")
            break;
        }
        error_check(err)

        line = strings.Trim(line, " ")

        // package name
        if strings.HasPrefix(line, "package ") {
            s := strings.Trim(line, "package")
            s  = strings.Trim(s, "\n")
            s  = strings.Trim(s, " ")
            s  = strings.Trim(s, ";")
            s  = strings.Trim(s, " ")
            p.name = s
        }
        // service
        if strings.HasPrefix(line, "service ") {
            _service := new(service)
            s := strings.Trim(line, "service")
            s  = strings.Trim(s, "\n")
            s  = strings.Trim(s, " ")
            s  = strings.Trim(s, "{")
            s  = strings.Trim(s, " ")
            _service.name = s
            p.ss = append(p.ss, _service)
        }
        // method 
        if strings.HasPrefix(line, "rpc ") {
            _service := p.ss[len(p.ss) - 1]
            _method := new(method)
            _method.index = len(_service.ms)

            s := strings.Trim(line, "rpc")
            s  = strings.Trim(s, " ")

            i  := 0
            ii := strings.Index(s, "(")
            _method.name = string(s[i:ii])

            // i  = strings.Index(s[ii+1:], "(")
            i  = ii + 1
            ii = i + strings.Index(s[i:], ")")
            fmt.Printf("i: %d, ii: %d\n", i, ii)
            _method.req_type_name = strings.Trim(s[i:ii], " ")

            i  = strings.Index(s[ii:], "(") + ii + 1
            ii = strings.Index(s[i:], ")") + i
            fmt.Printf("i: %d, ii: %d\n", i, ii)
            _method.rsp_type_name = strings.Trim(s[i:ii], " ")
            _service.ms = append(_service.ms, _method)
        }
    }

    handler(p)
}
