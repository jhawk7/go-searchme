node {
  checkout scm

  stage("build/push") {
    docker.withRegistry('https://hub.docker.com/', 'dockerhub') {
      def image = docker.build("jhawk7/go-searchme:${env.BUILD_ID}")
      image.push()
      image.push('latest')
    }
  }
}