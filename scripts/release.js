const request = require('request');
const fs = require('fs');
const base = 'http://zlepper.dk:3215';

/**
 * Gets an auth token for the given user
 * @param username The username
 * @param password The password
 */
exports.login = function (username, password) {
    return new Promise((resolve, reject) => {
        request.get(`${base}/api/auth/login?username=${username}&password=${password}`, (error, response, body) => {
            if (error) {
                return reject(error);
            }

            resolve(JSON.parse(body).token);
        })
    });
};

exports.createNewVersion = function ({token, name, notes = '', channel = 'stable'}) {
    return new Promise((resolve, reject) => {
        request.post({
            url: `${base}/api/version`,
            auth: {bearer: token},
            body: {name, notes, channel: {name: channel}},
            json: true
        }, (error, response, body) => {
            if (error) {
                return reject(error);
            }
            console.log(body);

            resolve({response: JSON.parse(body), token});
        })
    });
};

exports.addAsset = function ({version, platform, fileLocation}) {
    return new Promise((resolve, reject) => {
        const formData = {
            version,
            platform,
            file: fs.createReadStream(fileLocation)
        };

        request.post({url: `${base}/api/asset`, formData}, (error, response, body) => {
            if (error) {
                return reject(error);
            }

            resolve(JSON.parse(body));
        });
    });
};
