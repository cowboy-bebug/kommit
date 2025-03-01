class Kommit < Formula
    desc "Therapeutic approach to git commits with AI-generated messages"
    homepage "https://github.com/cowboy-bebug/kommitment"
    version "0.1.0"
    license "MIT"

    on_macos do
        if Hardware::CPU.arm?
            url "https://github.com/cowboy-bebug/kommitment/releases/download/v#{version}/kommit_v#{version}_darwin_arm64.tar.gz"
            sha256 "SHA256_DARWIN_ARM64"
        else
            url "https://github.com/cowboy-bebug/kommitment/releases/download/v#{version}/kommit_v#{version}_darwin_amd64.tar.gz"
            sha256 "SHA256_DARWIN_AMD64"
        end
    end

    on_linux do
      if Hardware::CPU.arm? && Hardware::CPU.is_64_bit?
        url "https://github.com/cowboy-bebug/kommitment/releases/download/v#{version}/kommit_v#{version}_linux_arm64.tar.gz"
        sha256 "SHA256_LINUX_ARM64"
      else
        url "https://github.com/cowboy-bebug/kommitment/releases/download/v#{version}/kommit_v#{version}_linux_amd64.tar.gz"
        sha256 "SHA256_LINUX_AMD64"
      end
    end

    depends_on "git"

    def install
      bin.install "kommit"
    end

    def caveats
      <<~EOS
        🧐 Your therapist is ready for session!

        To begin your repository's healing journey:
          kommit init

        Remember to set your OpenAI API key:
          export OPENAI_API_KEY=your_openai_api_key

        Or use a dedicated key for Kommitment:
          export KOMMIT_API_KEY=your_kommit_specific_key
      EOS
    end

    test do
      system "#{bin}/kommit", "version"
    end
  end
