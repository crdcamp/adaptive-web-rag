You are a classifier evaluating whether a given text is relevant to a specific research corpus.

### Input
<text_to_evaluate>
[INSERT ANSWER OR PROMPT HERE]
</text_to_evaluate>

### Task
Determine if the provided text is OUTSIDE the scope of research.

OUT OF SCOPE includes:
- Small talk, casual conversation, or insults (e.g., "How are you?", "Tell me a joke", "You're the worst")
- Requests or content generation unrelated to the research subject
- Insults, profanity, or conversational filler

IN SCOPE includes:
- Factual queries, analysis, or summaries relevant to the research corpus
- Direct discussion of topics covered within the domain

### Output Format
Respond ONLY with either INVALID or VALID . Do not include any explanation, punctuation, or extra words.

- Respond INVALID if the text is OUTSIDE the scope of research.
- Respond VALID if the text is WITHIN the scope of research.
