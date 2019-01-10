package rules_test

import (
	"io/ioutil"
	"reflect"
	"strings"
	"testing"

	"github.com/yoheimuta/protolint/internal/osutil"

	"github.com/yoheimuta/go-protoparser/parser/meta"

	"github.com/yoheimuta/protolint/internal/linter/file"

	"github.com/yoheimuta/protolint/internal/setting_test"

	"github.com/yoheimuta/protolint/internal/addon/rules"
	"github.com/yoheimuta/protolint/internal/linter/report"
)

func TestIndentRule_Apply(t *testing.T) {
	space4 := strings.Repeat(" ", 4)

	tests := []struct {
		name           string
		inputStyle     string
		inputProtoPath string
		wantFailures   []report.Failure
		wantExistErr   bool
	}{
		{
			name:           "correct syntax",
			inputProtoPath: setting_test.TestDataPath("rules", "indentrule", "syntax.proto"),
		},
		{
			name:           "incorrect syntax",
			inputStyle:     space4,
			inputProtoPath: setting_test.TestDataPath("rules", "indentrule", "incorrect_syntax.proto"),
			wantFailures: []report.Failure{
				report.Failuref(
					meta.Position{
						Filename: setting_test.TestDataPath("rules", "indentrule", "incorrect_syntax.proto"),
						Offset:   14,
						Line:     2,
						Column:   5,
					},
					`Found an incorrect indentation style "%s". "%s" is correct.`,
					space4,
					"",
				),
			},
		},
		{
			name:           "correct enum",
			inputProtoPath: setting_test.TestDataPath("rules", "indentrule", "enum.proto"),
		},
		{
			name:           "incorrect enum",
			inputProtoPath: setting_test.TestDataPath("rules", "indentrule", "incorrect_enum.proto"),
			wantFailures: []report.Failure{
				report.Failuref(
					meta.Position{
						Filename: setting_test.TestDataPath("rules", "indentrule", "incorrect_enum.proto"),
						Offset:   166,
						Line:     7,
						Column:   2,
					},
					`Found an incorrect indentation style "%s". "%s" is correct.`,
					" ",
					"",
				),
				report.Failuref(
					meta.Position{
						Filename: setting_test.TestDataPath("rules", "indentrule", "incorrect_enum.proto"),
						Offset:   69,
						Line:     4,
						Column:   9,
					},
					`Found an incorrect indentation style "%s". "%s" is correct.`,
					"        ",
					space4,
				),
				report.Failuref(
					meta.Position{
						Filename: setting_test.TestDataPath("rules", "indentrule", "incorrect_enum.proto"),
						Offset:   118,
						Line:     6,
						Column:   6,
					},
					`Found an incorrect indentation style "%s". "%s" is correct.`,
					"     ",
					space4,
				),
			},
		},
		{
			name:           "correct message",
			inputProtoPath: setting_test.TestDataPath("rules", "indentrule", "message.proto"),
		},
		{
			name:           "incorrect message",
			inputProtoPath: setting_test.TestDataPath("rules", "indentrule", "incorrect_message.proto"),
			wantFailures: []report.Failure{
				report.Failuref(
					meta.Position{
						Filename: setting_test.TestDataPath("rules", "indentrule", "incorrect_message.proto"),
						Offset:   106,
						Line:     6,
						Column:   5,
					},
					`Found an incorrect indentation style "%s". "%s" is correct.`,
					space4,
					strings.Repeat(space4, 2),
				),
				report.Failuref(
					meta.Position{
						Filename: setting_test.TestDataPath("rules", "indentrule", "incorrect_message.proto"),
						Offset:   166,
						Line:     9,
						Column:   1,
					},
					`Found an incorrect indentation style "%s". "%s" is correct.`,
					"",
					space4,
				),
				report.Failuref(
					meta.Position{
						Filename: setting_test.TestDataPath("rules", "indentrule", "incorrect_message.proto"),
						Offset:   311,
						Line:     14,
						Column:   13,
					},
					`Found an incorrect indentation style "%s". "%s" is correct.`,
					strings.Repeat(space4, 3),
					strings.Repeat(space4, 2),
				),
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			rule := rules.NewIndentRule(
				test.inputStyle,
				"\n",
				false,
			)

			proto, err := file.NewProtoFile(test.inputProtoPath, test.inputProtoPath).Parse()
			if err != nil {
				t.Errorf(err.Error())
				return
			}

			got, err := rule.Apply(proto)
			if test.wantExistErr {
				if err == nil {
					t.Errorf("got err nil, but want err")
				}
				return
			}
			if err != nil {
				t.Errorf("got err %v, but want nil", err)
				return
			}

			if !reflect.DeepEqual(got, test.wantFailures) {
				t.Errorf("got %v, but want %v", got, test.wantFailures)
			}
		})
	}
}

type testData struct {
	filePath   string
	originData []byte
}

func newTestData(
	fileName string,
) (testData, error) {
	filePath := setting_test.TestDataPath("rules", "indentrule", fileName)
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return testData{}, nil
	}
	return testData{
		filePath:   filePath,
		originData: data,
	}, nil
}

func (d testData) data() ([]byte, error) {
	return ioutil.ReadFile(d.filePath)
}

func (d testData) restore() error {
	newlineChar := "\n"
	lines := strings.Split(string(d.originData), newlineChar)
	return osutil.WriteLinesToExistingFile(d.filePath, lines, newlineChar)
}

func TestIndentRule_Apply_fix(t *testing.T) {
	space4 := strings.Repeat(" ", 4)

	correctSyntaxPath, err := newTestData("syntax.proto")
	if err != nil {
		t.Errorf("got err %v", err)
		return
	}

	incorrectSyntaxPath, err := newTestData("incorrect_syntax.proto")
	if err != nil {
		t.Errorf("got err %v", err)
		return
	}

	correctEnumPath, err := newTestData("enum.proto")
	if err != nil {
		t.Errorf("got err %v", err)
		return
	}

	incorrectEnumPath, err := newTestData("incorrect_enum.proto")
	if err != nil {
		t.Errorf("got err %v", err)
		return
	}

	correctMessagePath, err := newTestData("message.proto")
	if err != nil {
		t.Errorf("got err %v", err)
		return
	}

	incorrectMessagePath, err := newTestData("incorrect_message.proto")
	if err != nil {
		t.Errorf("got err %v", err)
		return
	}

	tests := []struct {
		name            string
		inputTestData   testData
		wantCorrectData testData
	}{
		{
			name:            "correct syntax",
			inputTestData:   correctSyntaxPath,
			wantCorrectData: correctSyntaxPath,
		},
		{
			name:            "incorrect syntax",
			inputTestData:   incorrectSyntaxPath,
			wantCorrectData: correctSyntaxPath,
		},
		{
			name:            "correct enum",
			inputTestData:   correctEnumPath,
			wantCorrectData: correctEnumPath,
		},
		{
			name:            "incorrect enum",
			inputTestData:   incorrectEnumPath,
			wantCorrectData: correctEnumPath,
		},
		{
			name:            "correct message",
			inputTestData:   correctMessagePath,
			wantCorrectData: correctMessagePath,
		},
		{
			name:            "incorrect message",
			inputTestData:   incorrectMessagePath,
			wantCorrectData: correctMessagePath,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			rule := rules.NewIndentRule(
				space4,
				"\n",
				true,
			)

			proto, err := file.NewProtoFile(test.inputTestData.filePath, test.inputTestData.filePath).Parse()
			if err != nil {
				t.Errorf(err.Error())
				return
			}

			_, err = rule.Apply(proto)
			if err != nil {
				t.Errorf("got err %v, but want nil", err)
				return
			}

			got, err := test.inputTestData.data()
			if !reflect.DeepEqual(got, test.wantCorrectData.originData) {
				t.Errorf(
					"got %s(%v), but want %s(%v)",
					string(got), got,
					string(test.wantCorrectData.originData), test.wantCorrectData.originData,
				)
			}

			err = test.inputTestData.restore()
			if err != nil {
				t.Errorf("got err %v", err)
			}
		})
	}
}
