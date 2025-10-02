"use client";

import { WIZARD_STEPS_ORDER, useWizardStore } from "@/store/wizard";

export function WizardStepper() {
  const { currentStepIndex, goTo } = useWizardStore();
  return (
    <div className="pt-0">
      <ol className="grid grid-cols-5 gap-2">
        {WIZARD_STEPS_ORDER.map((stepId, index) => {
          const completed = index < currentStepIndex;
          const isCurrent = index === currentStepIndex;
          return (
            <li key={stepId} className="flex flex-col items-center text-xs">
              <button
                type="button"
                className="flex flex-col items-center"
                onClick={() => completed && goTo(index)}
                aria-current={isCurrent ? "step" : undefined}
              >
                <span
                  className={`z-10 inline-flex h-7 w-7 items-center justify-center rounded-full border text-xs ${
                    isCurrent
                      ? "bg-secondary/70 text-foreground border-border"
                      : completed
                      ? "bg-muted text-foreground/80 border-border"
                      : "bg-card text-muted-foreground border-border"
                  }`}
                >
                  {index + 1}
                </span>
                <span className={`mt-2 capitalize ${isCurrent ? "text-foreground" : completed ? "text-foreground/70" : "text-muted-foreground"}`}>
                  {stepId}
                </span>
              </button>
            </li>
          );
        })}
      </ol>
    </div>
  );
}


