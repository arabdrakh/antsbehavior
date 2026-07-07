import subprocess
import sys
import os
import tempfile

_SRC_DIR = os.path.dirname(os.path.abspath(__file__))
_GO_DIR  = os.path.join(_SRC_DIR, "go")

_EXAMPLE00 = """\
4
##start
0 0 3
2 2 5
3 4 0
##end
1 8 3
0-2
2-3
3-1
"""


def _build_binary(output_path: str) -> bool:
    result = subprocess.run(
        ["go", "build", "-o", output_path, "."],
        cwd=_GO_DIR,
        capture_output=True,
        text=True,
    )
    if result.returncode != 0:
        print("=== go build failed ===", file=sys.stderr)
        print("stdout:", result.stdout, file=sys.stderr)
        print("stderr:", result.stderr, file=sys.stderr)
        print("========================", file=sys.stderr)
    return result.returncode == 0


def run_lemin(input_file: str) -> int:
    input_file = os.path.abspath(input_file)

    tmp_bin = tempfile.NamedTemporaryFile(delete=False, suffix=".exe" if os.name == "nt" else "")
    tmp_bin.close()
    bin_path = tmp_bin.name
    try:
        if not _build_binary(bin_path):
            return 1
        result = subprocess.run([bin_path, input_file])
        return result.returncode
    finally:
        try:
            os.unlink(bin_path)
        except OSError:
            pass


def main():
    print("Educational Practice Project")

    if len(sys.argv) > 1:
        input_file = sys.argv[1]
        if not os.path.exists(input_file):
            print(f"ERROR: file not found: {input_file}", file=sys.stderr)
            sys.exit(1)
        sys.exit(run_lemin(input_file))

    tmp_input = tempfile.NamedTemporaryFile(
        mode="w", suffix=".txt", delete=False
    )
    tmp_input.write(_EXAMPLE00)
    tmp_input.close()
    try:
        run_lemin(tmp_input.name)
    finally:
        try:
            os.unlink(tmp_input.name)
        except OSError:
            pass


if __name__ == "__main__":
    main()