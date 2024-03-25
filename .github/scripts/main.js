const https = require('https')
const {SQSClient, SendMessageCommand} = require('@aws-sdk/client-sqs')
const {DynamoDBClient, QueryCommand} = require('@aws-sdk/client-dynamodb')
const {unmarshall} = require('@aws-sdk/util-dynamodb')

const AWS_REGION = process.env.DEPLOYER_PLATFORM_AWS_REGION
const SQS_URL = process.env.DEPLOYER_PLATFORM_SQS_URL
const DYNAMO_TABLE = process.env.DEPLOYER_PLATFORM_DYNAMO_TABLE
const AWS_CREDS = {
    accessKeyId: process.env.DEPLOYER_PLATFORM_AWS_ACCESS_KEY_ID,
    secretAccessKey: process.env.DEPLOYER_PLATFORM_AWS_SECRET_ACCESS_KEY
};

const sqs = new SQSClient({region: AWS_REGION, credentials: AWS_CREDS})
const dynamodb = new DynamoDBClient({region: AWS_REGION, credentials: AWS_CREDS})


function queryForDeploymentStatus(messageId) {
    const query_params = {
        TableName: DYNAMO_TABLE,
        KeyConditionExpression: 'id = :id',
        FilterExpression: 'completed = :completed',
        ExpressionAttributeNames: {
            '#id': 'id',
            '#completed': 'completed',
            '#status': 'status',
            '#message': 'message',
        },
        ExpressionAttributeValues: {
            ':id': {
                S: messageId,
            },
            ':completed': {
                BOOL: true,
            },
        },
        ProjectionExpression: '#id, #completed, #status, #message',
        ScanIndexForward: false,  //returns items by descending timestamp
    }
    return new QueryCommand(query_params)
}

async function isDeploymentSuccessful(deploymentId, retries, waitSeconds) {
    for (let i = 0; i < retries; i++) {
        console.log(`Deployment pending, sleeping ${waitSeconds} seconds...`)
        await sleep(waitSeconds * 1000)

        try {
            const response = await dynamodb.send(queryForDeploymentStatus(deploymentId))
            console.log(`Query succeeded. Items found: ${response.Items.length}`)

            for (let i = 0; i < response.Items.length; i++) {
                const item = unmarshall(response.Items[i])
                if (item.completed) {
                    console.log(`Completed: ${item.id} - ${item.message} - ${item.completed} - ${item.status}`)
                    if (item.status === 'FAILED') {
                        console.error(`::error:: Deployment failed: ${item.message}`)
                        return false
                    }

                    return true
                }
            }
        } catch (err) {
            console.log(`Error querying table: ${err}`)
        }
    }
    return false
}

function sleep(ms) {
    return new Promise((resolve) => setTimeout(resolve, ms))
}

function main() {
    const url = process.env.TEST_DEFINITION_URL

    https.get(url, (res) => {
        let body = ''

        res.on('data', (chunk) => {
            body += chunk
        })

        res.on('end', async () => {
            let messageId
            try {
                const command = new SendMessageCommand({
                    QueueUrl: SQS_URL,
                    MessageBody: body,
                })
                data = await sqs.send(command)
                messageId = data.MessageId
                console.log(`Message sent: ${messageId}`)
            } catch (err) {
                console.error(`Error sending message: ${err}`)
            }

            // Execute the query with retries/sleeps
            let RETRIES = 200, WAIT_SECONDS = 15
            const success = await isDeploymentSuccessful(messageId, RETRIES, WAIT_SECONDS)
            if (!success) {
                process.exit(1)
            }
        })

        res.on('error', (err) => {
            console.error(`Error calling URL: ${err}`)
        })
    })
}

if (require.main === module) {
    main()
}
