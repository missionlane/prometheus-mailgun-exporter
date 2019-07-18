def group     = 'infra'
def service   = 'prometheus-mailgun-exporter'
def namespace = "${group}/${service}"

builderNode {
  stage("Build Docker Image") {
    checkout scm
    def build_date   = sh (script: 'date +%Y%m%d-%H:%M:%S', returnStdout: true).trim()
    def build_user   = sh (script: 'whoami', returnStdout: true).trim()
    def git_branch   = sh (script: 'git rev-parse --abbrev-ref HEAD', returnStdout: true).trim()
    def git_revision = sh (script: 'git rev-parse HEAD', returnStdout: true).trim()
    version          = sh (script: 'cat VERSION', returnStdout: true).trim()
    imageName = buildDockerImage(
      repository:  namespace,
      buildArgs: [
          "BUILD_DATE=${build_date}",
          "BUILD_USER=${build_user}",
          "GIT_BRANCH=${git_branch}",
          "GIT_REVISION=${git_revision}",
          'GO111MODULE=on',
          "VERSION=${version}"
      ]
    )
  }

  stage("Docker Promote All Builds") {
    promoteDockerImage(
      imageName: imageName,
      toTags: ["latest"]
    )
  }

  if (env.BRANCH_NAME == "master") {
    stage("Docker Promote Tag") {
      if (env.TAG_NAME ==~ /^v.+$/) {
        promoteDockerImage(
          imageName: imageName,
          toTags: ["latest", version]
        )
      }
    }
  }
}
