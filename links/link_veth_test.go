package links_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/srl-labs/containerlab/links"
	"github.com/srl-labs/containerlab/mocks/mocklinknodes"
	"github.com/srl-labs/containerlab/nodes/state"
)

func TestLinkVEthRaw_ToLinkBriefRaw(t *testing.T) {
	type fields struct {
		LinkCommonParams links.LinkCommonParams
		Endpoints        []*links.EndpointRaw
	}
	tests := []struct {
		name   string
		fields fields
		want   *links.LinkBriefRaw
	}{
		{
			name: "test1",
			fields: fields{
				LinkCommonParams: links.LinkCommonParams{
					MTU:    1500,
					Labels: map[string]string{"foo": "bar"},
					Vars:   map[string]any{"foo": "bar"},
				},
				Endpoints: []*links.EndpointRaw{
					{
						Node:  "node1",
						Iface: "eth1",
					},
					{
						Node:  "node2",
						Iface: "eth2",
					},
				},
			},
			want: &links.LinkBriefRaw{
				Endpoints: []string{"node1:eth1", "node2:eth2"},
				LinkCommonParams: links.LinkCommonParams{
					MTU:    1500,
					Labels: map[string]string{"foo": "bar"},
					Vars:   map[string]any{"foo": "bar"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &links.LinkVEthRaw{
				LinkCommonParams: tt.fields.LinkCommonParams,
				Endpoints:        tt.fields.Endpoints,
			}

			got := r.ToLinkBriefRaw()

			if d := cmp.Diff(got, tt.want); d != "" {
				t.Errorf("LinkVEthRaw.ToLinkBriefRaw() = %s", d)
			}
		})
	}
}

func TestLinkVEthRaw_GetType(t *testing.T) {
	type fields struct {
		LinkCommonParams links.LinkCommonParams
		Endpoints        []*links.EndpointRaw
	}
	tests := []struct {
		name   string
		fields fields
		want   links.LinkType
	}{
		{
			name: "test1",
			fields: fields{
				LinkCommonParams: links.LinkCommonParams{},
				Endpoints:        []*links.EndpointRaw{},
			},
			want: links.LinkTypeVEth,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &links.LinkVEthRaw{
				LinkCommonParams: tt.fields.LinkCommonParams,
				Endpoints:        tt.fields.Endpoints,
			}
			if got := r.GetType(); got != tt.want {
				t.Errorf("LinkVEthRaw.GetType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLinkVEthRaw_Resolve(t *testing.T) {
	// init Runtime Mock
	ctrl := gomock.NewController(t)

	// instantiate Mock Node 1
	fn1 := mocklinknodes.NewMockNode(ctrl)
	fn1.EXPECT().GetShortName().Return("node1").AnyTimes()
	var ept links.LinkEndpointType = links.LinkEndpointTypeVeth
	fn1.EXPECT().GetLinkEndpointType().Return(ept).AnyTimes()
	fn1.EXPECT().GetState().Return(state.Deployed).AnyTimes()

	// instantiate Mock Node 2
	fn2 := mocklinknodes.NewMockNode(ctrl)
	fn2.EXPECT().GetShortName().Return("node2").AnyTimes()
	fn2.EXPECT().GetLinkEndpointType().Return(ept).AnyTimes()
	fn2.EXPECT().GetState().Return(state.Deployed).AnyTimes()

	type fields struct {
		LinkCommonParams links.LinkCommonParams
		Endpoints        []*links.EndpointRaw
	}
	type args struct {
		params *links.ResolveParams
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *links.LinkVEth
		wantErr bool
	}{
		{
			name: "test1",
			fields: fields{
				LinkCommonParams: links.LinkCommonParams{
					MTU:    1500,
					Labels: map[string]string{"foo": "bar"},
					Vars:   map[string]any{"foo": "bar"},
				},
				Endpoints: []*links.EndpointRaw{
					{
						Node:  "node1",
						Iface: "eth1",
					},
					{
						Node:  "node2",
						Iface: "eth2",
					},
				},
			},
			args: args{
				params: &links.ResolveParams{
					Nodes: map[string]links.Node{
						"node1": fn1,
						"node2": fn2,
					},
				},
			},
			want: &links.LinkVEth{
				LinkCommonParams: links.LinkCommonParams{
					MTU:    1500,
					Labels: map[string]string{"foo": "bar"},
					Vars:   map[string]any{"foo": "bar"},
				},
				Endpoints: []links.Endpoint{
					&links.EndpointVeth{
						EndpointGeneric: links.EndpointGeneric{
							Node:      fn1,
							IfaceName: "eth1",
						},
					},
					&links.EndpointVeth{
						EndpointGeneric: links.EndpointGeneric{
							Node:      fn2,
							IfaceName: "eth2",
						},
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &links.LinkVEthRaw{
				LinkCommonParams: tt.fields.LinkCommonParams,
				Endpoints:        tt.fields.Endpoints,
			}
			got, err := r.Resolve(tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("LinkVEthRaw.Resolve() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			l := got.(*links.LinkVEth)
			if d := cmp.Diff(l.LinkCommonParams, tt.want.LinkCommonParams); d != "" {
				t.Errorf("LinkVEthRaw.Resolve() LinkCommonParams diff = %s", d)
			}

			for i, e := range l.Endpoints {
				if e.(*links.EndpointVeth).IfaceName != tt.want.Endpoints[i].(*links.EndpointVeth).IfaceName {
					t.Errorf("LinkVEthRaw.Resolve() EndpointVeth got %s, want %s", e.(*links.EndpointVeth).IfaceName, tt.want.Endpoints[i].(*links.EndpointVeth).IfaceName)
				}

				if e.(*links.EndpointVeth).Node != tt.want.Endpoints[i].(*links.EndpointVeth).Node {
					t.Errorf("LinkVEthRaw.Resolve() EndpointVeth got %s, want %s", e.(*links.EndpointVeth).Node, tt.want.Endpoints[i].(*links.EndpointVeth).Node)
				}
			}
		})
	}
}
