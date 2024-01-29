import os
import subprocess

def main():
    while True:
        print("Please select target OS:")
        print("1. Windows")
        print("2. Linux")
        os_choice = input("Enter your choice (1 or 2): ")

        if os_choice == "1":
            goos = "windows"
            goarch = "amd64"
            ext = ".exe"
            break
        elif os_choice == "2":
            goos = "linux"
            goarch = "amd64"
            ext = ""
            break
        else:
            print("Invalid choice. Please enter 1 or 2.")

    print(f"Compiling for {goos}...")
    os.environ["GOOS"] = goos
    os.environ["GOARCH"] = goarch
    output_filename = f"GoIris-{goos}{ext}"
    subprocess.run(["go", "build", "-ldflags=-s -w", "-trimpath", "-o", output_filename])
    print("Done.")

if __name__ == "__main__":
    main()
