Operating Mode: Workflow Verification
Workflow Files: @workflow/design.md, @workflow/prompt.md, @workflow/status.md, @workflow/followup.md
Task: Determine if the workflow results are ready to commit. You MUST FOLLOW the verification instructions below.

## Verification Instructions

Workflow verification is a READ-ONLY operation. Your purpose is to VERIFY workflow commit readiness, not to complete the implementation.

### Verification Rules

### REQUIRED BEHAVIORS
- ALL items in the verification checklist must pass
- ALWAYS LEVERAGE sub-agents to maximize efficiency

### Verification Checklist
- [ ] the project compiles successfully
- [ ] the integration tests pass
- [ ] all modified source code files are less than 1000 lines
- [ ] all implementation tasks MUST BE verified as complete

#### If ALL verification checklist items pass
- [ ] create new folder: workflow/completed/{git branch name}/{current UNIX timestamp}-{task name}
    - [ ] move workflow/design.md into new folder
    - [ ] move workflow/prompt.md into new folder
    - [ ] move workflow/status.md into new folder
    - [ ] move workflow/followup.md into new folder

#### If ANY verification checklist items fail
- [ ] replace workflow/followup.md with current failures and followup tasks