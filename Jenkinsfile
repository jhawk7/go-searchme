pipeline {
  agent {
    kubernetes {
      yaml '''
      apiVersion: v1
      kind: Pod
      metadata:
        name: dind
        namespace: devops
      spec:
        volumes:
          - name: docker-build-cache
            persistentVolumeClaim: 
              claimName: docker-build-cache
        nodeSelector:
          kubernetes.io/hostname: node7
        containers:
        - name: docker
          image: docker:latest
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
      steps {
        container('docker') {
          sh "docker build  -t jhawk7/go-searchme:$BUILD_NUMBER ."
          sh "docker push jhawk7/go-searchme:$BUILD_NUMBER"
          sh "docker push jhawk7/go-searchme:latest"
        }
      }
    }
  }
}
