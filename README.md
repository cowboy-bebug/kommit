# Kommit

> _Because your repository has commitment issues!_

Kommit is a therapeutic approach to git commits, helping your codebase express
itself through meaningful,
[conventional commit messages](https://www.conventionalcommits.org/en/v1.0.0/)
generated with AI.

## 🤔 Why Kommit?

Let's face it - most of us struggle with commitment issues when it comes to
writing
[good git commit messages](https://github.com/torvalds/subsurface-for-dirk/blob/724527bc3b660a9d54aab8e4dff50430450f1643/README.md?plain=1#L128-L150).
We've all been there:

```bash
git commit -m "fixed stuff"
git commit -m "updates"
git commit -m "it works now"
```

And who hasn't stared at that blinking cursor after `git commit -m` for what
feels like an eternity, mind suddenly blank, unable to summarize what you just
spent hours working on?

Your future self (and your teammates) deserve better. Kommit helps you write
clear, meaningful commit messages that tell the real story of your changes.

## 💊 Features

- **AI-Powered Analysis**: Analyzes your staged changes and generates
  meaningful, conventional commit messages
- **Relationship Counseling**: Improves communication between you and future
  developers
- **Emotional Intelligence**: Helps your code express what it's really trying to
  do
- **Commitment Structure**: Follows conventional commit formats for better git
  history
- **Scope Analysis**: Suggests relevant scopes based on your project structure
- **Self-Growth Options**: Choose to accept the AI suggestion, request another,
  or write your own

## 📦 Installation

### Using Homebrew (macOS and Linux)

```bash
brew install cowboy-bebug/tap/git-kommit
```

### From Source

Prerequisites:

- Go (1.24 or later required)
- Make

```bash
git clone https://github.com/cowboy-bebug/kommit.git
cd kommit
make install
```

This installs the `kommit` command-line tool.

### API Key Setup

Kommit requires an OpenAI API key to provide its therapeutic services:

```bash
# Set your OpenAI API key
export OPENAI_API_KEY=your_openai_api_key

# Alternatively, you can use a dedicated key for Kommit
export KOMMIT_API_KEY=your_kommit_specific_key
```

> **🔐 Therapy Privacy:** If both environment variables are set,
> `KOMMIT_API_KEY` takes precedence over `OPENAI_API_KEY`. This allows you to
> use a separate API key for Kommit if you prefer to keep your therapy sessions
> isolated from other OpenAI usage.

## 😌 Getting Started

### Initial Therapy Session

Begin your repository's healing journey:

```bash
git kommit init
```

This creates a `.kommitrc.yaml` configuration file. It'll analyze your project
structure and suggest meaningful scopes so your commits can finally express
themselves properly.

### Commit Therapy

When you're ready to commit changes:

1. Stage your changes with `git add` - **focus on logically grouping related
   changes**
2. Start a therapy session:

   ```bash
   git kommit
   ```

3. Review the AI-generated commit message
4. Choose to:
   - Accept the message and commit
   - Work through your commitment issues yourself (edit the message)
   - Request another therapy session for a better message (re-run)
   - Terminate the therapy session (exit without committing)

> **💡 Therapy Tip:** The key to successful commit therapy is to stage logically
> related changes together. Instead of worrying about writing the perfect commit
> message, focus on what belongs together in a commit. Let Kommit handle
> translating your changes into meaningful messages - that's what your therapist
> is here for!

## 🔍 How It Works

Kommit uses OpenAI's models to analyze your staged changes and generate
meaningful commit messages following conventional commit formats. It's like
couples therapy between you and your future self - improving how you communicate
today so there's less confusion tomorrow.

The therapy process:

1. Analyzes your staged changes
2. Considers your project's context and scopes
3. Generates a conventional commit message
4. Gives you options for proceeding
5. Handles the git commit process

## 🛠️ Configuration

After running `git kommit init`, you'll have a `.kommitrc.yaml` file with your
therapy plan:

```yaml
llm:
  model: gpt-4o-mini # Your therapist's qualifications
commit:
  types:
    - feat
    - fix
    - docs
    # ... other conventional commit types
  scopes:
    - ui
    - api
    - auth
    # ... project-specific scopes
```

## 💭 Examples

**Before therapy:**

```bash
$ git commit -m "fixed the thing that was broken yesterday"
$ git log -1 --format=medium
commit a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6q7r8s9t0
Author: Developer Name <dev@example.com>
Date:   Mon Mar 4 10:30:45 2023 -0800

    fixed the thing that was broken yesterday
```

**After therapy**:

```bash
$ git kommit
$ git log -1 --format=medium
commit a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6q7r8s9t0
Author: Developer Name <dev@example.com>
Date:   Mon Mar 4 10:30:45 2023 -0800

    feat(auth): implement JWT token refresh mechanism

    - Added token refresh endpoint to authentication service
    - Implemented automatic refresh when token is near expiration
    - Updated documentation with new token lifecycle

    [Generated by Kommit]
```

## 🙏 Contributing

Your contributions to this therapeutic journey are welcome! Whether you're
fixing bugs, adding features, or improving documentation, we appreciate your
help in making git history a better place.

## 📜 License

MIT - Because everyone deserves access to commit message therapy.

---

_Remember: Good commit messages are like good apologies - specific, sincere, and
they don't include the phrase "various changes". Your future self will thank you
for the therapy._
