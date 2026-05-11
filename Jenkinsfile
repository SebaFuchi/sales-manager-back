pipeline {
    agent any
    environment {
        PORT = '8080'
        DB = credentials('CONNECTION_URL')
        FIREBASE_CREDENTIALS = credentials('FIREBASE_SERVICE_ACCOUNT_JSON') // JSON secret string in Jenkins
        DOCKER_TAG = 'sales-manager-back'
        EXTERNAL_PORT = '8083'
        DB_NAME = 'database_name'
    }
    stages {
        stage('Build') {
            steps {
                echo 'Building..'
                sh 'docker build -t $DOCKER_TAG .'
            }
        }
        stage('Deploy') {
            steps {
                echo 'Stopping previous version...'
                sh 'docker stop $DOCKER_TAG || true'
                sh 'docker rm $DOCKER_TAG || true'
                echo 'Deploying....'
                sh 'docker run -d -e DB -e PORT -e FIREBASE_CREDENTIALS_JSON="$FIREBASE_CREDENTIALS" -p $EXTERNAL_PORT:8080 --name $DOCKER_TAG $DOCKER_TAG'
            }
        }
    }
}
