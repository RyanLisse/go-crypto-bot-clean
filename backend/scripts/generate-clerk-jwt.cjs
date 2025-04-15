const https = require('https');

const data = JSON.stringify({
  user_id: 'user_2vUFyhJiiREYf2qe6hTY0UiY9Kw'
});

const options = {
  hostname: 'api.clerk.com',
  path: '/v1/sessions',
  method: 'POST',
  headers: {
    'Authorization': 'Bearer sk_test_lYeQgfVtLJ9shBkRnK7e5Uhlt4SbJD7IX2fYILUfP1',
    'Content-Type': 'application/json',
    'Content-Length': data.length
  }
};

const req = https.request(options, res => {
  let body = '';
  res.on('data', d => body += d);
  res.on('end', () => {
    const json = JSON.parse(body);
    console.log('JWT:', json.jwt);
  });
});

req.on('error', error => {
  console.error(error);
});

req.write(data);
req.end();