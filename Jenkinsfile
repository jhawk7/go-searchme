pipeline {
  agent {
    kubernetes {
      yaml '''
      apiVersion: v1
      kind: Pod
      metadata:
        name: kaniko
        namespace: devops
      spec:
        nodeSelector:
          kubernetes.io/hostname: node7
        containers:
        - name: kaniko
          image: gcr.io/kaniko-project/executor:debug
          imagePullPolicy: Always
          command:
          - sleep
          args:
          - 9999999
          volumeMounts:
            - name: docker-config
              mountPath: /kaniko/.docker
            - name: docker-cache
              mountPath: /kaniko/
          envFrom:
          - secretRef:
              name: github-secret
        restartPolicy: Never
        volumes:
          - name: docker-config
          - name: docker-cache
            persistentVolumeClaim:
              claimName: docker-build-cache
      '''
    }
  }

  stages {
    stage("build/deploy") {
      steps {
        container(name: 'kaniko', shell: '/busybox/sh') {
          sh '''#!/busybox/sh
          /kaniko/executor --dockerfile=Dockerfile \
          --context=git://github.com/jhawk7/go-searchme.git \
          --cache=true
          --custom-platform=linux/arm64 \
          --destination=jhawk7/go-searchme:$BUILD_ID \
          --destination=jhawk7/go-searchme:latest'''
        }
      }
    }
  }
}
