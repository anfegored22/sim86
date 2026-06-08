# Agent Instructions

The user is working to improve memory for programming, especially Go.

Primary goal: help the user become able to build Go programs without internet access or searching.

Rules:

- When a question about programming or Go appears, always add one or more quiz entries to `programming_qa.md`.
- When the user writes `TRACK: <question>`, always add one or more quiz entries for that question to `programming_qa.md`.
- Quiz entries should be short, concrete, and memory-oriented, not open-ended explanations.
- Prefer multiple-choice questions with one correct answer.
- Do not simply copy the user's wording if it is broad or open-ended. Convert it into quiz-style questions that test specific facts or procedures.
- It is acceptable and encouraged to create multiple quiz questions from one user question when several facts are worth memorizing.
- If the user asks a programming question and code is changed or commands are run, still record the quiz entries.

Quiz format:

```markdown
## YYYY-MM-DD

### Q: <short quiz question>

a) <option>
b) <option>
c) <option>
d) <option>

Answer: <letter>) <option>
```
