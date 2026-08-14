class DonkCli < Formula
  desc "DONK is a keyboard-first terminal AI workspace"
  homepage "https://donk-cli.com"
  url "https://github.com/richavery/donk-cli/archive/refs/tags/v#{version}.tar.gz"
  version "1.1.4"
  license "FSL-1.1-MIT"

  depends_on "go" => :build

  def install
    ldflags = "-s -w -X github.com/richavery/donk-cli/internal/version.Version=#{version}"
    system "go", "build", *std_go_args(ldflags:, output: bin/"donk-cli"), "."

    generate_completions_from_executable(bin/"donk-cli", "completion")
    man1.install Utils.safe_popen_read(bin/"donk-cli", "man").pipe("gzip -c"), "#{name}.1.gz"
  end

  test do
    assert_match "DONK", shell_output("#{bin}/donk-cli --version")
  end
end
