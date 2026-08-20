---
description: Start work on a milestone with its definition of done in context
argument-hint: "<M0|M1|M2|M3|M4|M5|M6>"
---

Work milestone $ARGUMENTS from `docs/MILESTONES.md`.

Before writing code:

1. Read the milestone entry and restate its definition of done in one sentence.
2. Read `CLAUDE.md` hard constraints and name any that this milestone touches.
3. List the files you plan to change. Do not create new top-level packages.
4. Check whether this milestone needs a change to `model/event.go`. If it does, stop and write a `docs/DECISIONS.md` entry first, then ask before proceeding.

While working:

- Replace the milestone-tagged TODOs in the stubs rather than writing new files alongside them.
- Populate `Confidence` on every observation you emit. An unpopulated `Confidence` is a bug.
- Write the test that proves the definition of done before the implementation.

When the definition of done passes, stop. Report what landed and what the next milestone needs. Do not continue into the next milestone.
