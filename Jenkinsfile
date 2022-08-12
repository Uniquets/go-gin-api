pipeline {
    agent any

    environment {
        PROJECT = 'go-demo'
        REPOSITROY_URL = 'uniquets'+'/godemo'
        VERSION = '1.0.0'
    }

    stages {
        stage('登录docker') {
            steps {
                withCredentials([usernamePassword(credentialsId: '79a0ae06-c977-4a5b-97a0-c26f4eb79f3f', passwordVariable: 'ts1170998607!', usernameVariable: 'uniquets')]) {
                    sh "docker login --username=${USERNAME} --password=${PASSWORD}"
                }
            }
        }

        stage('打包镜像') {
            steps {
                   sh 'docker build -t ${REPOSITROY_URL}:${VERSION} .'
            }
        }

        stage('上传镜像') {
            steps {
                sh 'docker push ${REPOSITROY_URL}:${VERSION}'
            }
        }
        stage('部署服务') {
            steps {
                sh runCommand() + "sudo docker stop ${PROJECT}-${params.profile} || true"
                sh runCommand() + "sudo docker rm ${PROJECT}-${params.profile} || true"
                sh runCommand() + "sudo docker rmi ${REPOSITROY_URL}:${VERSION} || true"
                sh runCommand() + dockerRun()
            }
        }
    }
}

String runCommand() {
    def command = "sshpass  ssh -p ${params.port} -o \"StrictHostKeyChecking no\" ${params.username}@${params.host} "
    echo "command: ${command}"
    return "${command}";
}

String dockerRun() {
    String dockerRun = "sudo docker run --name ${PROJECT}-${params.profile} -p ${params.expose_port}:${params.expose_port} -d --network host ${REPOSITROY_URL}:${VERSION} "
    return dockerRun;
}