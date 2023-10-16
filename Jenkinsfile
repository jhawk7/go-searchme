pipeline {
  agent {
    kubernetes {
      yaml '''
      apiVersion: v1
      kind: Pod
      spec:
        volumes:
          - name: docker-build-cache
            persistentVolumeClaim: 
              claimName: docker-build-cache
        serviceAccountName: jenkins-agents
        containers:
        - name: docker
          image: docker/docker:latest
          volumeMounts:
          - name: docker-build-cache
            mountPath: /var/lib/docker
            subPath: docker
          command:
          - cat
          tty: true
          securityContext:
            privileged: true
      '''
    }
  }
  stages {
    stage('build/push') {
      container('docker') {
        sh "docker build  -t jhawk7/go-searchme:$BUILD_NUMBER ."
        sh "docker push jhawk7/go-searchme:$BUILD_NUMBER"
      }
    }
  }
}
