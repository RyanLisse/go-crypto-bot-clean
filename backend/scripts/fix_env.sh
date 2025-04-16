#!/bin/bash

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}=== Checking and fixing .env file ===${NC}"

# Check if .env file exists
if [ ! -f ".env" ]; then
    echo -e "${RED}Error: .env file not found!${NC}"
    exit 1
fi

# Create a temporary file
TEMP_FILE=$(mktemp)

# Process the .env file line by line
while IFS= read -r line; do
    # Skip empty lines and comments
    if [[ -z "$line" || "$line" =~ ^\# ]]; then
        echo "$line" >> "$TEMP_FILE"
        continue
    fi

    # Extract key and value
    KEY=$(echo "$line" | cut -d '=' -f 1)
    VALUE=$(echo "$line" | cut -d '=' -f 2-)

    # If this is a MEXC credential, trim any spaces
    if [[ "$KEY" == "MEXC_API_KEY" || "$KEY" == "MEXC_SECRET_KEY" || "$KEY" == "MEXC_CRED_ENCRYPTION_KEY" ]]; then
        # Trim spaces
        TRIMMED_VALUE=$(echo "$VALUE" | xargs)
        
        # Check if the value has changed
        if [ "$VALUE" != "$TRIMMED_VALUE" ]; then
            echo -e "${YELLOW}Fixed spaces in $KEY${NC}"
            echo "$KEY=$TRIMMED_VALUE" >> "$TEMP_FILE"
        else
            echo "$line" >> "$TEMP_FILE"
        fi
    else
        # Keep the line as is
        echo "$line" >> "$TEMP_FILE"
    fi
done < ".env"

# Replace the original file
mv "$TEMP_FILE" ".env"

echo -e "${GREEN}Fixed .env file. Please try running the application again.${NC}"

# Run a basic check on the API key to detect any obvious issues
API_KEY=$(grep "^MEXC_API_KEY=" .env | cut -d '=' -f 2 | xargs)
if [[ $API_KEY == *" "* ]]; then
    echo -e "${RED}Warning: MEXC_API_KEY still contains spaces after trimming!${NC}"
fi

if [[ $API_KEY == *":"* || $API_KEY == *";"* || $API_KEY == *","* ]]; then
    echo -e "${RED}Warning: MEXC_API_KEY contains invalid characters (colon, semicolon, or comma)!${NC}"
fi

# Print the API key length to check if it's reasonable
echo -e "${YELLOW}MEXC_API_KEY length: ${#API_KEY} characters${NC}"
echo -e "${YELLOW}First 5 characters: ${API_KEY:0:5}...${NC}"
echo -e "${YELLOW}Last 4 characters: ...${API_KEY: -4}${NC}"

if [[ ${#API_KEY} -lt 10 || ${#API_KEY} -gt 100 ]]; then
    echo -e "${RED}Warning: MEXC_API_KEY length (${#API_KEY}) seems unusual!${NC}"
fi

# Check encryption key
ENCRYPTION_KEY=$(grep "^MEXC_CRED_ENCRYPTION_KEY=" .env | cut -d '=' -f 2 | xargs)
echo -e "${YELLOW}MEXC_CRED_ENCRYPTION_KEY length: ${#ENCRYPTION_KEY} characters${NC}"

# Suggest next steps
echo -e "${GREEN}===== Next steps =====${NC}"
echo -e "1. Run the server using: ${YELLOW}./run_server.sh${NC}"
echo -e "2. If issues persist, run the test client: ${YELLOW}go run cmd/test_mexc_client/main.go${NC}" 