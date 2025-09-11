package jobs

import "testing"

func TestReconcileEnvironmentStatus_NoChangeWhenSameAfterNormalization(t *testing.T) {
	changed, target := reconcileEnvironmentStatus("active", "running")
	if changed {
		t.Fatalf("expected no change, got change to %q", target)
	}
}

func TestReconcileEnvironmentStatus_ChangeWhenDifferent(t *testing.T) {
	changed, target := reconcileEnvironmentStatus("creating", "ready")
	if !changed || target != "active" {
		t.Fatalf("expected change to active, got changed=%v target=%q", changed, target)
	}
}
