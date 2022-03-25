package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"

	"sigs.k8s.io/kustomize/kyaml/yaml"
)

const (
	configMapKind = iota + 1
	secretKind
)

var headerMap = map[int]string{
	configMapKind: `
apiVersion: v1
kind: ConfigMap
`,
	secretKind: `
apiVersion: v1
kind: Secret
type: Opaque
`,
}

type ConfigMap struct {
	name string
}

func (cm *ConfigMap) read(r io.Reader) (map[string]string, error) {
	return readDataMap(r)
}

func (cm *ConfigMap) write(w io.Writer, m map[string]string) error {
	return buildWrite(cm.name, configMapKind)(w, m)
}

type Secret struct {
	name string
}

func (s *Secret) read(r io.Reader) (map[string]string, error) {
	dm, err := readDataMap(r)
	if err != nil {
		return nil, err
	}
	return decodeBase64Values(dm)
}

func (sw *Secret) write(w io.Writer, m map[string]string) error {
	return buildWrite(sw.name, secretKind)(w, m)
}

func buildWrite(name string, kind int) func(io.Writer, map[string]string) error {
	return func(w io.Writer, m map[string]string) (err error) {
		var (
			header string
			ok     bool
		)
		if header, ok = headerMap[kind]; !ok {
			panic(fmt.Sprintf("unknown kind value: %d", kind))
		}

		rn, err := yaml.Parse(header)
		if err != nil {
			return err
		}
		if _, err := rn.Pipe(yaml.SetK8sName(name)); err != nil {
			return err
		}

		switch kind {
		case configMapKind:
			err = rn.LoadMapIntoConfigMapData(m)
		case secretKind:
			err = rn.LoadMapIntoSecretData(m)
		default:
			panic(fmt.Sprintf("unknown kind value: %d", kind))
		}
		if err != nil {
			return err
		}

		s, err := rn.String()
		if err != nil {
			return err
		}
		_, err = fmt.Fprint(w, s)
		return
	}
}

func decodeBase64Values(dm map[string]string) (map[string]string, error) {
	m := make(map[string]string, len(dm))
	for k, v := range dm {
		b, err := base64.StdEncoding.DecodeString(v)
		if err != nil {
			return nil, err
		}
		m[k] = string(b)
	}
	return m, nil
}

func readDataMap(r io.Reader) (map[string]string, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	rn, err := yaml.Parse(string(b))
	if err != nil {
		return nil, err
	}
	return rn.GetDataMap(), nil
}
