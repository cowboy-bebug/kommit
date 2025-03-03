class Kommit < Formula
  desc "Therapeutic approach to git commits with AI-generated messages"
  homepage "https://github.com/cowboy-bebug/kommit"
  license "MIT"

  head "https://github.com/cowboy-bebug/kommit.git", branch: "main"

  stable do
    url "https://github.com/cowboy-bebug/kommit.git",
        tag:      "v0.1.0",
        revision: "52b1845a42476310f7e5ab84eb56e090b5c7268a"
  end

  depends_on "go" => :build
  depends_on "git"

  def install
    system "go", "build",
           "-ldflags", "-X main.version=#{version} -X main.commit=#{stable.specs[:revision][0,7]} -X main.date=#{Time.now.utc.strftime("%Y-%m-%dT%H:%M:%SZ")}",
           "-o", "git-kommit"

    bin.install "git-kommit"
  end

  def caveats
    <<~EOS
      🧐 Your therapist is ready for session!

      To begin your repository's healing journey:
        git kommit init

      Remember to set your OpenAI API key:
        export OPENAI_API_KEY=your_openai_api_key

      Or use a dedicated key for Kommit:
        export KOMMIT_API_KEY=your_kommit_specific_key
    EOS
  end

  test do
    system "#{bin}/git", "kommit", "version"
  end
end
