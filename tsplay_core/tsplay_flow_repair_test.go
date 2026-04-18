package tsplay_core

import "testing"

func TestBuildFlowRepairContextUsesNestedFailurePath(t *testing.T) {
	contextPayload, err := BuildFlowRepairContext(FlowRepairContextOptions{
		Flow: &Flow{
			SchemaVersion: CurrentFlowSchemaVersion,
			Name:          "nested_failure",
			Steps: []FlowStep{
				{
					Action: "if",
					Condition: &FlowStep{
						Action:   "assert_visible",
						Selector: "#panel",
					},
					Then: []FlowStep{
						{
							Name:     "open export menu",
							Action:   "click",
							Selector: "#menu",
						},
						{
							Name:     "click export",
							Action:   "click",
							Selector: `text="Old export"`,
						},
					},
				},
			},
		},
		Result: &FlowResult{
			Vars: map[string]any{
				"orders_url": "https://example.com/orders",
			},
			Trace: []FlowStepTrace{
				{
					Index:  1,
					Path:   "1",
					Action: "if",
					Status: "error",
					Condition: &FlowStepTrace{
						Path:   "1.condition",
						Action: "assert_visible",
						Status: "ok",
					},
					Children: []FlowStepTrace{
						{
							Path:   "1.then.1",
							Action: "click",
							Status: "ok",
						},
						{
							Path:        "1.then.2",
							Name:        "click export",
							Action:      "click",
							Status:      "error",
							Args:        map[string]any{"selector": `text="Old export"`},
							ArgsSummary: `{"selector":"text=\"Old export\""}`,
							Error:       "locator click: timeout",
						},
					},
				},
			},
		},
	})
	if err != nil {
		t.Fatalf("build repair context: %v", err)
	}
	if contextPayload.FailedStepPath != "1.then.2" {
		t.Fatalf("failed_step_path = %q", contextPayload.FailedStepPath)
	}
	if contextPayload.FailedStep == nil {
		t.Fatalf("expected failed step")
	}
	if contextPayload.FailedStep.Path != "1.then.2" {
		t.Fatalf("failed step path = %q", contextPayload.FailedStep.Path)
	}
	if contextPayload.FailedStep.Step.Name != "click export" {
		t.Fatalf("failed step = %#v", contextPayload.FailedStep.Step)
	}
	if len(contextPayload.NearbySteps) != 2 {
		t.Fatalf("nearby steps = %#v", contextPayload.NearbySteps)
	}
	if contextPayload.NearbySteps[0].Path != "1.then.1" || contextPayload.NearbySteps[0].Relation != "previous" {
		t.Fatalf("unexpected previous nearby step: %#v", contextPayload.NearbySteps[0])
	}
	if contextPayload.NearbySteps[1].Path != "1.then.2" || contextPayload.NearbySteps[1].Relation != "failed" {
		t.Fatalf("unexpected failed nearby step: %#v", contextPayload.NearbySteps[1])
	}
}
