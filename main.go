package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/open-policy-agent/frameworks/constraint/pkg/client"
	"github.com/open-policy-agent/frameworks/constraint/pkg/client/drivers/local"
	"github.com/open-policy-agent/frameworks/constraint/pkg/core/templates"
	"github.com/open-policy-agent/gatekeeper/pkg/target"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	yamlutil "k8s.io/apimachinery/pkg/util/yaml"
)

func main() {
	//fmt.Println(string(y))
	ctx := context.Background()
	client, err := createClient()

	t := ReadConstraintTemplate()

	if _, err = client.AddTemplate(ctx, t); err != nil {
		fmt.Println(err)
		os.Exit(4)
	}

	c := ReadConstraint()
	if _, err = client.AddConstraint(ctx, c); err != nil {
		fmt.Println(err)
		os.Exit(7)
	}

	d := ReadData()
	if _, err = client.AddData(ctx, d); err != nil {
		fmt.Println(err)
		os.Exit(9)
	}

	if _, err = client.AddData(ctx, t); err != nil {
		fmt.Println(err)
		os.Exit(9)
	}

	if _, err = client.AddData(ctx, c); err != nil {
		fmt.Println(err)
		os.Exit(9)
	}

	resps, err := client.Audit(ctx)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("\nOpa Audit response\n", resps)

	if len(resps.Results()) > 0 {
		fmt.Println("violations found!")
		for k, v := range resps.ByTarget {
			fmt.Println(k, v.Input, v.Target, v.Trace, v.Results)
		}
		os.Exit(1)
	}

	fmt.Println("good to go!")
	os.Exit(0)
}

func createClient() (*client.Client, error) {
	target := &target.K8sValidationTarget{}
	driver := local.New()
	backend, err := client.NewBackend(client.Driver(driver))
	if err != nil {
		return nil, err
	}
	c, err := backend.NewClient(client.Targets(target))
	if err != nil {
		return nil, err
	}
	return c, nil
}

func ReadConstraintTemplate() *templates.ConstraintTemplate {
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
	obj := &templates.ConstraintTemplate{}
	if err := json.Unmarshal(rawObj.Raw, obj); err != nil {
		fmt.Println(err)
		os.Exit(3)
	}
	fmt.Println("\nConstraint Template:\n", obj)
	return obj
}

func ReadConstraint() *unstructured.Unstructured {
	y, err := ioutil.ReadFile("opatemplates/constraint.yaml")
	if err != nil {
		fmt.Println(err)
		os.Exit(5)
	}

	decoder := yamlutil.NewYAMLToJSONDecoder(bytes.NewReader(y))

	// read a document from our multidoc yaml file
	var rawObj runtime.RawExtension
	if err := decoder.Decode(&rawObj); err != nil {
		fmt.Println(err)
		os.Exit(6)
	}

	// decode using unstructured JSON scheme
	obj := &unstructured.Unstructured{}
	if err := json.Unmarshal(rawObj.Raw, obj); err != nil {
		fmt.Println(err)
		os.Exit(7)
	}
	fmt.Println("\nConstraint:\n", obj)
	//fmt.Println(obj.GetAPIVersion())
	return obj
}

func ReadData() *unstructured.Unstructured {
	y, err := ioutil.ReadFile("opatemplates/data.yaml")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	decoder := yamlutil.NewYAMLToJSONDecoder(bytes.NewReader(y))

	// read a document from our multidoc yaml file
	var rawObj runtime.RawExtension
	if err := decoder.Decode(&rawObj); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// decode using unstructured JSON scheme
	obj := &unstructured.Unstructured{}
	if err := json.Unmarshal(rawObj.Raw, obj); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("\nSample Data:\n", obj)
	return obj
}
