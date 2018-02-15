Release checklist:
 * [ ] Write release notes in NEWS. Get them reviewed and merged.
 * [ ] Ensure your local copy is up to date with master.
 * [ ] Ensure you can sign commits and that any yubikeys/smartcards you
       need for that are plugged in.
 * [ ] Run `./release.sh vX.Y.Z` which will:
   * Verify the version you specified is a valid semver
   * Ensure your working directory is clean
   * Run tests
   * Tag the release and sign the tag
   * Verify the the tag
   * Build release artifacts
 * [ ] Push the tag to Github with:
 ```
 git push vX.Y.Z origin
 ```
 * [ ] Sign the release artifacts by running
 ```
 gpg -u gpg --local-user 0xCDDE268EBB729EC7 --detach-sign --armor <path to artifact>
 ```
   for each release artifact. Do _not_ try to sign all of them at once by globbing. If you do, gpg will
   sign the combination of all the release artifacts instead of each one individually.
 * [ ] Create a draft release on Github for the tag and upload all the release artifacts and
       their signatures. Copy and paste the release notes from NEWS here as well.
 * [ ] Get someone to review the draft release
 * [ ] Publish!
 * [ ] Tag the docker image that was automatically built when the NEWS PR was merged
   * [ ] Run `git describe` to determine which docker container to tag.
   * [ ] Run `docker pull quay.io/coreos/ct:<git-descibe-output>`
   * [ ] Run `docker tag quay.io/coreos/ct:<git-describe-output> quay.io/coreos/ct:vX.Y.Z`
   * [ ] Run `docker push quay.io/coreos/ct:vX.Y.Z`
