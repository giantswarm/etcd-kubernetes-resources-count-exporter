package collector

import "testing"

func Test_parseLine(t *testing.T) {
	tests := []struct {
		name          string
		line          string
		wantNamespace string
		wantKind      string
		wantErr       bool
	}{
		{
			name:          "Pod",
			line:          "pods/akv2k8s/akv2k8s-controller-799f4db8c4-g4ltn",
			wantNamespace: "akv2k8s",
			wantKind:      "pods",
			wantErr:       false,
		},
		{
			name:          "Persistent Volume",
			line:          "persistentvolumes/pvc-18a19af7-d1ac-40f1-b153-f3a8c24f33f0",
			wantNamespace: "Not namespaced",
			wantKind:      "persistentvolumes",
			wantErr:       false,
		},
		{
			name:          "Cluster Role Bindings",
			line:          "clusterrolebindings/akv2k8s-controller",
			wantNamespace: "Not namespaced",
			wantKind:      "clusterrolebindings",
			wantErr:       false,
		},
		{
			name:          "Cluster Issuers",
			line:          "cert-manager.io/clusterissuers/letsencrypt-giantswarm",
			wantNamespace: "Not namespaced",
			wantKind:      "clusterissuers.cert-manager.io",
			wantErr:       false,
		},
		{
			name:          "Empty line",
			line:          "",
			wantNamespace: "",
			wantKind:      "",
			wantErr:       true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotNamespace, gotKind, err := parseLine(tt.line)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseLine() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotNamespace != tt.wantNamespace {
				t.Errorf("parseLine() gotNamespace = %v, want %v", gotNamespace, tt.wantNamespace)
			}
			if gotKind != tt.wantKind {
				t.Errorf("parseLine() gotKind = %v, want %v", gotKind, tt.wantKind)
			}
		})
	}
}
