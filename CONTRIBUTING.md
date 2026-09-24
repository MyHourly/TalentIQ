# Contributing to TalentIQ

## 1. Choose your assigned module branch

Use only the branch assigned to your module unless coordinated otherwise.

Examples:

```text
feature/talent-profile
feature/certificate-and-assessment
feature/skills
feature/skill-graph
feature/ai-semantic-talent-search
feature/skill-gap-analysis
feature/technology-trends
feature/talent-recommendations
feature/dashboard
```

## 2. Start from the latest branch

```bash
git checkout feature/<your-module>
git pull origin feature/<your-module>
```

## 3. Work locally

Implement your changes, run the appropriate formatter/tests/build, and keep your changes within the module where possible.

## 4. Review before committing

```bash
git status
git diff
```

## 5. Commit clearly

```bash
git add .
git commit -m "feat: <short description>"
```

## 6. Push your branch

```bash
git push origin feature/<your-module>
```

## 7. Create a Pull Request

Open a Pull Request from your feature branch into `develop`.

PR descriptions should include:

- What changed
- Why it changed
- How it was tested
- Any API/database/event contract changes
- Any follow-up work still required

## 8. Before merging

Make sure:

- Build succeeds
- Tests pass
- No secrets are committed
- No unrelated module changes are included
- API/event changes are documented
- Reviewer feedback is addressed

## 9. Integration flow

```text
feature branch -> Pull Request -> develop -> integration testing -> reviewed PR -> main
```

Do not push normal feature work directly to `main`.

## 10. Collaboration

Frontend and backend members working on the same module should agree on API contracts before integration. For event-driven changes, document Kafka topic/event names, RabbitMQ queues, payload shape, producer/consumer ownership and failure handling.
