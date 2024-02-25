// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package extractrequest defines an Analyzer that serves as a trivial
// example and test of the Analysis API. It reports a diagnostic for
// every call to a function or method of the name specified by its
// -name flag. It also exports a fact for each declaration that
// matches the name, plus a package-level fact if the package contained
// one or more such declarations.
package extractrequest

import (
	"autoarmor/armorgo/analysisutil"
	"encoding/json"
	"fmt"
	"go/ast"
	"go/types"
	"golang.org/x/tools/go/analysis"
	"io/ioutil"
	"strings"
)

var projectName = "callerservice"

// var projectName = "productcatalogservice"
//var projectName = "checkoutservice"

// var projectName = "shippingservice"
//var projectName = "catalogue"

//var projectName = "user"

//var projectName = "payment"

const Doc = `find modeled gRPCs and TCP invocations in the go project

The extractrequest analysis reports calls to the modeled functions.`

const jsonPath = "/usr/local/go/src/autogo/armorgo/resources/rpc_info.json"

var Analyzer = &analysis.Analyzer{
	Name:             "extractrequest",
	Doc:              Doc,
	Run:              run,
	RunDespiteErrors: true,
	FactTypes:        []analysis.Fact{new(foundFact)},
}

type Manifest struct {
	Service  string    `json:"service"`
	Version  string    `json:"version"`
	Requests []Request `json:"requests"`
}

type Request struct {
	Type string `json:"type"`
	URL  string `json:"url"`
	Name string `json:"name"`
	Path string `json:"path"`
}

type TCPManifest struct {
	Service  string       `json:"service"`
	Version  string       `json:"version"`
	Requests []TCPRequest `json:"requests"`
}

type TCPRequest struct {
	Type string `json:"type"`
	URL  string `json:"url"`
	Name string `json:"name"`
	Port string `json:"port"`
}

type endpoint struct {
	name string
	path string
}

type tcpEndpoint struct {
	//pkgName string
	name string
	//path    string
	port string
}

var tcpMethods = map[string]struct{ args, results []string }{
	// "Flush": {{}, {"error"}}, // http.Flusher and jpeg.writer conflict
	"Ping": {[]string{}, []string{"error"}},                              //sqlx.Ping, mgo.Ping
	"Open": {[]string{"string", "string"}, []string{"*sql.DB", "error"}}, // sqlx.Open
	//"Close": {[]string{}, []string{"error"}},                                      // sqlx.Close

	"Select":  {[]string{"interface{}", "string", "[]interface{}"}, []string{"error"}}, // sqlx.Select
	"Get":     {[]string{"interface{}", "string", "[]interface{}"}, []string{"error"}}, // sqlx.Get
	"Query":   {[]string{"string", "[]interface{}"}, []string{"*sql.Rows", "error"}},   // sqlx.Query
	"Prepare": {[]string{"string"}, []string{"*sql.Stmt", "error"}},                    // sqlx.Prepare

	"UpsertId":    {[]string{"interface{}", "interface{}"}, []string{"*mgo.ChangeInfo", "error"}}, // mgo.session.UpsertId
	"RemoveAll":   {[]string{"interface{}"}, []string{"*mgo.ChangeInfo", "error"}},                // mgo.RemoveAll
	"Update":      {[]string{"interface{}", "interface{}"}, []string{"error"}},                    // mgo.Update
	"Find":        {[]string{"interface{}"}, []string{"*mgo.Query"}},                              // mgo.Find
	"FindId":      {[]string{"interface{}"}, []string{"*mgo.Query"}},                              // mgo.FindId
	"UpdateAll":   {[]string{"interface{}", "interface{}"}, []string{"*mgo.ChangeInfo", "error"}}, // mgo.UpdateAll
	"Remove":      {[]string{"interface{}"}, []string{"error"}},                                   // mgo.Remove
	"EnsureIndex": {[]string{"mgo.Index"}, []string{"error"}},                                     // mgo.EnsureIndex

	//"MarshalJSON":   {[]string{}, []string{"[]byte", "error"}},                            // json.Marshaler
	//"MarshalXML":    {[]string{"*xml.Encoder", "xml.StartElement"}, []string{"error"}},    // xml.Marshaler
	//"ReadByte":      {[]string{}, []string{"byte", "error"}},                              // io.ByteReader
	//"ReadFrom":      {[]string{"=io.Reader"}, []string{"int64", "error"}},                 // io.ReaderFrom
	//"ReadRune":      {[]string{}, []string{"rune", "int", "error"}},                       // io.RuneReader
	//"Scan":          {[]string{"=fmt.ScanState", "rune"}, []string{"error"}},              // fmt.Scanner
	//"Seek":          {[]string{"=int64", "int"}, []string{"int64", "error"}},              // io.Seeker
	//"UnmarshalJSON": {[]string{"[]byte"}, []string{"error"}},                              // json.Unmarshaler
	//"UnmarshalXML":  {[]string{"*xml.Decoder", "xml.StartElement"}, []string{"error"}},    // xml.Unmarshaler
	//"UnreadByte":    {[]string{}, []string{"error"}},
	//"UnreadRune":    {[]string{}, []string{"error"}},
	//"WriteByte":     {[]string{"byte"}, []string{"error"}},                // jpeg.writer (matching bufio.Writer)
	//"WriteTo":       {[]string{"=io.Writer"}, []string{"int64", "error"}}, // io.WriterTo
}

var servicePaths = map[string]tcpEndpoint{
	"catalogue": tcpEndpoint{"catalogue-db", "3306"},
	"user":      tcpEndpoint{"user-db", "27017"},
}

var name string // -name flag
var deploymentFilePrefix = "/usr/local/go/src/autogo/armorgo/resources/deployment_files/example/"
var outputPrefix = "output/"
var model map[string]*endpoint
var tcpModel map[string]*tcpEndpoint
var ManifestData Manifest
var TCPManifestData TCPManifest
var flag = false

func init() {
	//model
	model = make(map[string]*endpoint)
	tcpModel = make(map[string]*tcpEndpoint)
	grpcs := ModelFromJson(jsonPath)
	fmt.Println("[AutoArmor]: grpc package - " + grpcs.Package)
	rpcs := grpcs.Rpcs
	for _, rpc := range rpcs {
		targetName := strings.ToLower(rpc.Service)
		path := rpc.Path
		//server := rpc.Service + "Server"
		//key := server + "." + rpc.Name
		key := rpc.Name
		model[key] = &endpoint{targetName, path}
	}

	//tcpModel["Open"] = &tcpEndpoint{"github.com/jmoiron/sqlx", "Open", "", ""}
	//tcpModel["Close"] = &tcpEndpoint{"github.com/jmoiron/sqlx", "Close", "", ""}
	//tcpModel["Ping"] = &tcpEndpoint{"github.com/jmoiron/sqlx", "Ping", "", ""}

	//parse deployment file

	deploymentFilePath := deploymentFilePrefix + projectName + ".yaml"
	deployment := ParseYaml(deploymentFilePath)
	version := deployment.Metadata.Labels.Version
	ManifestData = Manifest{Version: version, Service: projectName}
	TCPManifestData = TCPManifest{Version: version, Service: projectName}
	//todo:combine
	fmt.Println("[AutoArmor]: analyzing project - " + projectName + ":" + version + "\n")

	Analyzer.Flags.StringVar(&name, "name", name, "name of the function to find")
}

func run(pass *analysis.Pass) (interface{}, error) {
	pass.IgnoredFiles = append(pass.IgnoredFiles, "server_test.go")
	//inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	// Fast path: if the package doesn't import net/http,
	// skip the traversal.
	if analysisutil.Imports(pass.Pkg, "testing") {
		return nil, nil
	}

	for _, f := range pass.Files {
		if strings.HasSuffix(f.Name.Name, "_test.go") {
			continue
		}
		ast.Inspect(f, func(n ast.Node) bool {
			if call, ok := n.(*ast.CallExpr); ok {
				var id *ast.Ident
				var x ast.Expr
				switch fun := call.Fun.(type) {
				case *ast.Ident:
					id = fun
				case *ast.SelectorExpr:
					id = fun.Sel
					x = fun.X
				}
				//name2 := "EmptyCart"
				if id != nil && x != nil && !pass.TypesInfo.Types[id].IsType() {
					//methodName := id.Name
					//var typeExpr ast.Expr
					//switch xx := x.(type) {
					//case *ast.TypeAssertExpr:
					//	typeExpr = xx.Type
					//}
					//var typ string
					//switch typeID := typeExpr.(type) {
					//case *ast.Ident:
					//	typ = typeID.Name
					//}
					//key2 := typ + "." + methodName
					//fmt.Println(key2)
					key := id.Name
					endpoint, found := model[key]
					_, foundTCP := tcpMethods[key]
					if found && isGrpc(pass.TypesInfo, call) {
						//fmt.Println("nice")
						request := Request{
							Type: "grpc",
							URL:  endpoint.name, // can add port, but ok if not
							Name: endpoint.name,
							Path: endpoint.path,
						}
						ManifestData.Requests = append(ManifestData.Requests, request)
						flag = true
						pass.Report(analysis.Diagnostic{
							Pos:     call.Lparen,
							Message: fmt.Sprintf("call of %s(...)", key),
							SuggestedFixes: []analysis.SuggestedFix{{
								Message: fmt.Sprintf("Add '_TEST_'"),
								TextEdits: []analysis.TextEdit{{
									Pos:     call.Lparen,
									End:     call.Lparen,
									NewText: []byte("_TEST_"),
								}},
							}},
						})
						//} else if foundTCP && key=="Find" {//testing point
					} else if foundTCP && isTcp(pass, pass.TypesInfo, call, key) { //testing point
						fmt.Println("[AutoArmor]: find call " + key)
						tcpInfo, _ := servicePaths[projectName]
						tcpRequest := TCPRequest{
							Type: "tcp",
							URL:  tcpInfo.name + ":" + tcpInfo.port,
							Name: tcpInfo.name,
							Port: tcpInfo.port,
						}
						if !flag {
							TCPManifestData.Requests = append(TCPManifestData.Requests, tcpRequest)
							flag = true
						}
						pass.Report(analysis.Diagnostic{
							Pos:     call.Lparen,
							Message: fmt.Sprintf("call of %s(...)", key),
							SuggestedFixes: []analysis.SuggestedFix{{
								Message: fmt.Sprintf("Add '_TEST_'"),
								TextEdits: []analysis.TextEdit{{
									Pos:     call.Lparen,
									End:     call.Lparen,
									NewText: []byte("_TEST_"),
								}},
							}},
						})
					}
				}

			}
			return true
		})
	}

	// Export a fact for each matching function.
	//
	// These facts are produced only to test the testing
	// infrastructure in the analysistest package.
	// They are not consumed by the extractrequest Analyzer
	// itself, as would happen in a more realistic example.
	for _, f := range pass.Files {
		for _, decl := range f.Decls {
			if decl, ok := decl.(*ast.FuncDecl); ok && decl.Name.Name == "AddItem" {
				if obj, ok := pass.TypesInfo.Defs[decl.Name].(*types.Func); ok {
					pass.ExportObjectFact(obj, new(foundFact))
				}
			}
		}
	}

	if len(pass.AllObjectFacts()) > 0 {
		pass.ExportPackageFact(new(foundFact))
	}

	if flag && ManifestData.Requests != nil {
		Output(ManifestData)
	}

	if flag && TCPManifestData.Requests != nil {
		TCPOutput(TCPManifestData)
		//flag = false
	}

	return nil, nil
}

// 类型判断，减少误报
func isGrpc(info *types.Info, expr *ast.CallExpr) bool {
	fun, _ := expr.Fun.(*ast.SelectorExpr)
	sig, _ := info.Types[fun].Type.(*types.Signature)
	if sig == nil {
		return false // the call is not of the form x.f()
	}

	params := sig.Params()
	if params.Len() != 3 {
		return false // the function called does not return three values.
	}
	ptr := params.At(2).Type()
	ppt := ptr.(*types.Slice)
	if ppt == nil || !isNamedType(ppt.Elem(), "google.golang.org/grpc", "CallOption") {
		return false // the first return type is not *http.Response.
	}
	return true
}

// 类型判断，减少误报
func isTcp(pass *analysis.Pass, info *types.Info, expr *ast.CallExpr, key string) bool {
	fun, _ := expr.Fun.(*ast.SelectorExpr)
	sig, _ := info.Types[fun].Type.(*types.Signature)
	if sig == nil {
		return false // the call is not of the form x.f()
	}
	args := sig.Params()
	results := sig.Results()
	expect, _ := tcpMethods[key]

	//fmt.Println("get")

	if matchParams(pass, expect.args, args, "") && matchParams(pass, expect.results, results, "") {
		return true
	}
	return false
}

func isNamedType(t types.Type, path, name string) bool {
	n, ok := t.(*types.Named)
	if !ok {
		return false
	}
	obj := n.Obj()
	return obj.Name() == name && obj.Pkg() != nil && obj.Pkg().Path() == path
}

func Output(mani Manifest) {
	b, err := json.MarshalIndent(mani, "", " ")
	if err != nil {
		fmt.Println("error:", err)
	}
	//生成json文件
	err = ioutil.WriteFile(outputPrefix+mani.Service+".json", b, 0777)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	var data interface{}
	err = json.Unmarshal(b, &data)
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Println("data", data)
}

func TCPOutput(mani TCPManifest) {
	b, err := json.MarshalIndent(mani, "", " ")
	if err != nil {
		fmt.Println("error:", err)
	}
	//生成json文件
	err = ioutil.WriteFile(outputPrefix+mani.Service+".json", b, 0777)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	var data interface{}
	err = json.Unmarshal(b, &data)
	if err != nil {
		fmt.Println("error:", err)
	}
	//fmt.Println("data", data)
}

// Does each type in expect with the given prefix match the corresponding type in actual?
func matchParams(pass *analysis.Pass, expect []string, actual *types.Tuple, prefix string) bool {
	for i, x := range expect {
		if !strings.HasPrefix(x, prefix) {
			continue
		}
		if i >= actual.Len() {
			return false
		}
		if !matchParamType(x, actual.At(i).Type()) {
			return false
		}
	}
	if prefix == "" && actual.Len() > len(expect) {
		return false
	}
	return true
}

// Does this one type match?
func matchParamType(expect string, actual types.Type) bool {
	expect = strings.TrimPrefix(expect, "=")
	// Overkill but easy.
	act := typeString(actual)
	return act == expect
}

func typeString(typ types.Type) string {
	return types.TypeString(typ, (*types.Package).Name)
}

// foundFact is a fact associated with functions that match -name.
// We use it to exercise the fact machinery in tests.
type foundFact struct{}

func (*foundFact) String() string {
	return "found"
}
func (*foundFact) AFact() {}
