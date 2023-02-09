package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type JsonData struct {
	Data map[string]interface{} `json:"-"`
}

func JsonDecode(r io.Reader) (JsonData, error) {
	d := json.NewDecoder(r)
	var jsData JsonData
	if err := d.Decode(&jsData.Data); err != nil {
		return JsonData{}, err
	}
	return jsData, nil
}

func Generate(data map[string]interface{}, location string) error {
	// fmt.Printf("%+v\n",data)
	componentData := make(map[string]string, 0)

	for key, value := range data {

		a, err := json.MarshalIndent(value, "", "      ")
		if err != nil {
			return err
		}
		componentData[key] = string(a)
	}
	// fmt.Printf("\n%+v", componentData)

	for k, val := range componentData {
		if _, err := os.Stat(filepath.Join(location, strings.ToLower(k))); os.IsNotExist(err) {
			// path/to/whatever does not exist
			if err := os.MkdirAll(filepath.Join(location, strings.ToLower(k)), os.ModePerm); err != nil {
				return err
			}
			//create formio js file
			err := formIoJs(k, filepath.Join(location, strings.ToLower(k)))
			if err != nil {
				return err
			}

			//create unit test js file
			err = unitTestJs(k, filepath.Join(location, strings.ToLower(k)))
			if err != nil {
				return err
			}

			//create unit test js file
			err = mainJs(k, filepath.Join(location, strings.ToLower(k)), val)
			if err != nil {
				return err
			}

			//create editForm folder and js file
			err = editForm(filepath.Join(location, strings.ToLower(k)), k)
			if err != nil {
				return err
			}

			//create fixtures folder and js file
			err = fixtures(filepath.Join(location, strings.ToLower(k)), val)
			if err != nil {
				return err
			}
		}
		break
	}
	return nil
}

func formIoJs(filename string, loc string) error {
	jsContent := `import ~filename~EditDisplay from './editForm/~filename~.edit.display';
export default function() {
  return {
    components: [
      {
        type: 'tabs',
        key: 'tabs',
        components: [
          {
            label: '~filename~',
            key: 'display',
            weight: 0,
            components: ~filename~EditDisplay
          }
        ]
      }
    ]
  };
}`
	jsContent = strings.ReplaceAll(jsContent, "~filename~", filename)
	err := os.WriteFile(filepath.Join(loc, filename+".form.js"), []byte(jsContent), 0777)
	if err != nil {
		return err
	}
	return nil
}

func unitTestJs(filename string, loc string) error {
	jsContent := `import ~filename~Component from './~filename~';

import {
  comp1
} from './fixtures';

describe('~filename~ Component', () => {
  it('Should build a ~filename~ component in builder mode', (done) => {
    new ~filename~Component(comp1, {
      builder: true
    });
    done();
  });
});`
	jsContent = strings.ReplaceAll(jsContent, "~filename~", filename)
	err := os.WriteFile(filepath.Join(loc, filename+".unit.js"), []byte(jsContent), 0777)
	if err != nil {
		return err
	}
	return nil
}

func mainJs(filename string, loc string, val string) error {
	jsContent := `import Component from '../_classes/component/Component';

export default class ~filename~Component extends Component {
  static schema() {
    return %v;
  }

  static get builderInfo() {
    return {
      title: '~filename~',
      icon: 'cubes',
      group: 'IPeG',
      weight: 120,
      schema: ~filename~Component.schema()
    };
  }

  get defaultSchema() {
    return ~filename~Component.schema();
  }
}`
	jsContent = strings.ReplaceAll(jsContent, "~filename~", filename)
	jsContent = fmt.Sprintf(jsContent, val)
	err := os.WriteFile(filepath.Join(loc, filename+".js"), []byte(jsContent), 0777)
	if err != nil {
		return err
	}
	return nil
}

func editForm(location string, filename string) error {
	if _, err := os.Stat(filepath.Join(location, "editForm")); os.IsNotExist(err) {

		if err := os.MkdirAll(filepath.Join(location, "editForm"), os.ModePerm); err != nil {
			return err
		}

		jsContent := `export default [
  {
    key: '~filename~ComponentDescription',
    label: '~filename~ component description',
    input: false,
    tag: 'p',
    content: '~filename~ components can be used to render special fields or widgets inside your app. ' +
      'For information on how to display in an app, see ' +
      '<a href="http://help.form.io/userguide/#custom" target="_blank">' +
      'custom component documentation' +
      '</a>.',
    type: 'htmlelement',
    weight: 5
  },
  {
    type: 'textarea',
    as: 'json',
    editor: 'ace',
    weight: 10,
    input: true,
    key: 'componentJson',
    label: 'Custom Element JSON',
    tooltip: 'Enter the JSON for this custom element.'
  }
];`
		jsContent = strings.ReplaceAll(jsContent, "~filename~", filename)
		err := os.WriteFile(filepath.Join(location, "editForm", filename+".edit.display.js"), []byte(jsContent), 0777)
		if err != nil {
			return err
		}
	}
	return nil
}

func fixtures(location string, val string) error {
	if _, err := os.Stat(filepath.Join(location, "fixtures")); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Join(location, "fixtures"), os.ModePerm); err != nil {
			return err
		}

		comp1 := fmt.Sprintf(`export default %v`, val)
		err := os.WriteFile(filepath.Join(location, "fixtures", "comp1.js"), []byte(comp1), 0777)
		if err != nil {
			return err
		}

		index := "export comp1 from './comp1';"
		err = os.WriteFile(filepath.Join(location, "fixtures", "index.js"), []byte(index), 0777)
		if err != nil {
			return err
		}
	}
	return nil
}
