# To learn more about how to use Nix to configure your environment
# see: https://firebase.google.com/docs/studio/customize-workspace
{ pkgs, ... }: {
  # Which nixpkgs channel to use.
  channel = "stable-24.05"; # or "unstable"

  # Use https://search.nixos.org/packages to find packages
  packages = [
    # Go and related tools
    pkgs.go_1_24
    pkgs.gcc
    pkgs.sqlite
    pkgs.air # Hot reloading for Go

    # Version control and utilities
    pkgs.git
    pkgs.curl
    pkgs.jq # JSON processing
    pkgs.ca-certificates
    pkgs.tzdata

    # Frontend development
    pkgs.nodejs_20
    pkgs.bun # Package manager used in frontend

    # Docker for containerization
    pkgs.docker
    pkgs.docker-compose
  ];

  # Sets environment variables in the workspace
  env = {
    # Go configuration
    GOPATH = "$HOME/go";

    # Application configuration
    CONFIG_PATH = "./backend/configs";
    CONFIG_FILE = "config.yaml";
    LOG_LEVEL = "debug";

    # Database configuration
    DB_PATH = "./backend/data/dev.db";
    DATABASE_ENABLED = "true";

    # MEXC API configuration
    MEXC_BASE_URL = "https://api.mexc.com";
    MEXC_WEBSOCKET_URL = "wss://wbs.mexc.com/ws";
  };

  # Enable services
  services = {
    # Enable Docker
    docker.enable = true;

    # Enable SQLite database
    postgres = {
      enable = true;
      extensions = ["pgvector"];
    };
  };

  idx = {
    # Search for the extensions you want on https://open-vsx.org/ and use "publisher.id"
    extensions = [
      "golang.go"
      "ms-azuretools.vscode-docker"
      "dbaeumer.vscode-eslint"
      "bradlc.vscode-tailwindcss"
      "esbenp.prettier-vscode"
      "rangav.vscode-thunder-client"
    ];

    # Enable previews
    previews = {
      enable = true;
      previews = {
        backend = {
          command = ["cd" "backend" "&&" "go" "run" "."];
          manager = "web";
          env = {
            PORT = "8080";
          };
        };
        frontend = {
          command = ["cd" "frontend" "&&" "bun" "run" "dev"];
          manager = "web";
          env = {
            PORT = "3000";
          };
        };
      };
    };

    # Workspace lifecycle hooks
    workspace = {
      # Runs when a workspace is first created
      onCreate = {
        install-backend-deps = "cd backend && go mod download";
        install-frontend-deps = "cd frontend && bun install";
      };
      # Runs when the workspace is (re)started
      onStart = {
        check-dirs = "mkdir -p backend/data/logs backend/configs";
      };
    };
  };
}
