package v1beta1

import (
	"testing"

	"github.com/grafana/loki/pkg/logql"
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
			wantErr: false,
			args: args{
				ns:   "foo",
				expr: `1`,
			},
		},
		{
			name:    "literal binary operators",
			wantErr: false,
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
