package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/GoogleContainerTools/kpt-functions-sdk/go/pkg/framework/types"
	"github.com/open-policy-agent/frameworks/constraint/pkg/core/templates"
	"k8s.io/apimachinery/pkg/runtime"
	yamlutil "k8s.io/apimachinery/pkg/util/yaml"
)

func main() {
	configs := types.Configs{}

	configs = append(configs, ReadConstraintTemplate())

	fmt.Println(configs)
}

func ReadConstraintTemplate() templates.ConstraintTemplate {
	y, err := ioutil.ReadFile("opatemplates/template.yaml")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	decoder := yamlutil.NewYAMLToJSONDecoder(bytes.NewReader(y))

	// read a document from our multidoc yaml file
	var rawObj runtime.RawExtension
	if err := decoder.Decode(&rawObj); err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	// decode using unstructured JSON scheme
	obj := templates.ConstraintTemplate{}
	if err := json.Unmarshal(rawObj.Raw, obj); err != nil {
		fmt.Println(err)
		os.Exit(3)
	}
	fmt.Println("\nConstraint Template:\n", obj)
	return obj
}
