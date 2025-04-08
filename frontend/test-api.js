// Simple script to test API connectivity
import fetch from 'node-fetch';

const API_URL = 'http://localhost:8080/api/v1/status';

async function testAPI() {
  try {
    console.log(`Testing API connection to: ${API_URL}`);
    const response = await fetch(API_URL);

    console.log(`Response status: ${response.status} ${response.statusText}`);

    if (response.ok) {
      const data = await response.json();
      console.log('Response data:', JSON.stringify(data, null, 2));
    } else {
      console.error('Failed to fetch data');
    }
  } catch (error) {
    console.error('Error:', error.message);
  }
}

testAPI();
