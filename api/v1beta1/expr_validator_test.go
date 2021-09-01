package v1beta1

import (
	"reflect"
	"testing"

	"github.com/grafana/loki/pkg/logql"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Test_enforceNode(t *testing.T) {
	type args struct {
		ns   string
		expr string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "literal",
			wantErr: true,
			args: args{
				ns:   "foo",
				expr: `1`,
			},
		},
		{
			name:    "literal binary operators",
			wantErr: true,
			args: args{
				ns:   "foo",
				expr: `1 + 1`,
			},
		},
		{
			name:    "log range",
			wantErr: false,
			args: args{
				ns:   "foo",
				expr: `absent_over_time({job="abcd"}[5m])`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node, err := logql.ParseExpr(tt.args.expr)
			if err != nil {
				t.Errorf("enforceNode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err := enforceNode(tt.args.ns, node); (err != nil) != tt.wantErr {
				t.Errorf("enforceNode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGlobalLokiRule_ValidateExpressions(t *testing.T) {
	type fields struct {
		TypeMeta   metav1.TypeMeta
		ObjectMeta metav1.ObjectMeta
		Spec       GlobalLokiRuleSpec
		Status     GlobalLokiRuleStatus
	}
	tests := []struct {
		name    string
		fields  fields
		want    *GlobalLokiRuleSpec
		wantErr bool
	}{
		{
			name: "validate expressions",
			fields: fields{
				Spec: GlobalLokiRuleSpec{
					Groups: []*LokiRuleGroup{
						{
							Name: "foo",
							Rules: []*LokiGroupRule{
								{
									Alert: "new-expresssion",
									Expr:  `absent_over_time({job="abcd"}[5m])`,
									For:   "5m",
								},
								{
									Alert: "literal",
									Expr:  `1`,
									For:   "5m",
								},
							},
						},
					},
				},
			},
			want: &GlobalLokiRuleSpec{
				Groups: []*LokiRuleGroup{
					{
						Name: "foo",
						Rules: []*LokiGroupRule{
							{
								Alert: "new-expresssion",
								Expr:  `absent_over_time({job="abcd"}[5m])`,
								For:   "5m",
							},
							{
								Alert: "literal",
								Expr:  `1`,
								For:   "5m",
							},
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lokiRule := &GlobalLokiRule{
				TypeMeta:   tt.fields.TypeMeta,
				ObjectMeta: tt.fields.ObjectMeta,
				Spec:       tt.fields.Spec,
				Status:     tt.fields.Status,
			}
			got, err := lokiRule.ValidateExpressions()
			if (err != nil) != tt.wantErr {
				t.Errorf("GlobalLokiRule.ValidateExpressions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GlobalLokiRule.ValidateExpressions() = %v, want %v", got, tt.want)
			}
		})
	}
}
