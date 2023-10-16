pipeline {
  agent {
    kubernetes {
      yaml '''
      apiVersion: v1
      kind: Pod
      metadata:
        name: kaniko
        namespace: kaniko
      spec:
        nodeSelector:
          kubernetes.io/hostname: node7
        containers:
        - name: kaniko-demo
          image: gcr.io/kaniko-project/executor:latest
          args: ["--context=git://github.com/jhawk7/go-searchme.git",
            "--destination=jhawk7/go-searchme: "${BUILD_NUMBER}"
            "--destination=jhawk7/go-searchme:latest",
            "--dockerfile=Dockerfile",
            "--context=dir://workspace",
            "--custom-platform=linux/arm64"]
          volumeMounts:
            - name: docker-config
              mountPath: /kaniko/.docker
            - name: docker-cache
              mountPath: /workspace
          envFrom:
          - secretRef:
              name: github-secret
        restartPolicy: Never
        volumes:
          - name: docker-config
          - name: dockerfile-storage
            persistentVolumeClaim:
              claimName: docker-build-cache
      '''
    }
  }
  stages {
    stage('build/push') {
      steps {
        sh "echo image jhawk7/go-searchme:$BUILD_NUMBER pushed"
      }
    }
  }
}
