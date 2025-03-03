class Kommit < Formula
    desc "Therapeutic approach to git commits with AI-generated messages"
    homepage "https://github.com/cowboy-bebug/kommit"
    version "0.1.0"
    license "MIT"

    on_macos do
        if Hardware::CPU.arm?
            url "https://github.com/cowboy-bebug/kommit/releases/download/v#{version}/kommit_v#{version}_darwin_arm64.tar.gz"
            sha256 "3d952e69607d935c4f61513f675ca3db9403ad8f9861f0483ce08b857c068157"
        else
            url "https://github.com/cowboy-bebug/kommit/releases/download/v#{version}/kommit_v#{version}_darwin_amd64.tar.gz"
            sha256 "8ac8b3562ec8a2738ce89ac1081a0c5ce091d5b5210140f264c803df3386e312"
        end
    end

    on_linux do
      if Hardware::CPU.arm? && Hardware::CPU.is_64_bit?
        url "https://github.com/cowboy-bebug/kommit/releases/download/v#{version}/kommit_v#{version}_linux_arm64.tar.gz"
        sha256 "6e49900a693bf8db3f29a0185a70405b7ab751e40ca2a57aae4f25738dbb8629"
      else
        url "https://github.com/cowboy-bebug/kommit/releases/download/v#{version}/kommit_v#{version}_linux_amd64.tar.gz"
        sha256 "b150358071f1897dbdf6844cbedec9928bfffc9f5e057a9e8351d82f89075657"
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

        Or use a dedicated key for Kommit:
          export KOMMIT_API_KEY=your_kommit_specific_key
      EOS
    end

    test do
      system "#{bin}/kommit", "version"
    end
  end
