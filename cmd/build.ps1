 convert the above shell script into a ps1 script

# Get the commit message from the passed in argument
$COMMIT_MESSAGE = $args[0]
$MAJOR = 0
$MINOR = 0

# if there is no commit message, set it to "no message"
if ($COMMIT_MESSAGE -eq $null) {
    $COMMIT_MESSAGE = "no message"
}

# create a function with a parameter called message
# that will be displayed in the terminal
function quit($message) {
    # print the message to the terminal
    Write-Host " 🛑 $message"
    exit 1
}

# create a function to do the git stuff
function gitstep {
    git add -A                          | Out-Null
    git commit -m "$COMMIT_MESSAGE"     | Out-Null
    git push                            | Out-Null
    # then checkout the Beor and throw away output to avoid printing it
    git checkout Beor                   | Out-Null
    git merge main                      | Out-Null

    git add -A                          | Out-Null
    git commit -m "$COMMIT_MESSAGE"     | Out-Null

    # git push origin beor
    git push origin Beor                | Out-Null
}

function returntoMain {
    # then checkout the main and throw away output to avoid printing it
    git checkout main                   | Out-Null
    git merge Beor                      | Out-Null
    git push origin main                | Out-Null
}


# create a function to do the version stuff
function versionstep {
	# check if there is a VERSION.yaml file if there is not then create one
	if (!(Test-Path -Path ".\cmd\VERSION.yaml")) {
		# create the VERSION.yaml file
		New-Item -Path ".\cmd\VERSION.yaml" -ItemType File -Force
		# set the major version to 0
		$MAJOR=0
		# set the minor version to 0
		$MINOR=0

	} else {
		# get the version numbers from the VERSION.yaml file as numbers
		$MAJOR=(Get-Content -Path ".\cmd\VERSION.yaml" | Select-String -Pattern ("major:") | ForEach-Object { $_.Line.Split(":")[1] }).Trim()
		$MINOR=(Get-Content -Path ".\cmd\VERSION.yaml" | Select-String -Pattern ("minor:") | ForEach-Object { $_.Line.Split(":")[1] }).Trim()

		# change MAJOR and MINOR to numbers
		$MAJOR = [int]$MAJOR
		$MINOR = [int]$MINOR

                # Check if -M flag is passed in
		if ($args[1] -eq "-M") {
		    # increment the major version
		    $MAJOR = $MAJOR + 1
		    # set the minor version to 0
		    $MINOR = 0
		} else {
		    # increment the minor version
		    $MINOR = $MINOR + 1
		}

		# write the new version numbers to the VERSION.yaml file
		$version = "major: $MAJOR`nminor: $MINOR"
		$version | Out-File -FilePath ".\cmd\VERSION.yaml"

		# print the version numbers to the terminal
		Write-Host " 📦 Version: $MAJOR.$MINOR"
	}
}

# create the stage messages
$STAGE0 = " 🌳 Checking out Beor branch                   "
$STAGE1 = " 🔁 Updating VERSION                           "
$STAGE2 = " 📦 Logging into ECR                           "
$STAGE3 = " 🐳 Building docker image                      "
$STAGE4 = " 📑 Tagging docker image                       "
$STAGE5 = " 📌 Pushing docker image to ECR                "
$STAGE6 = " 🏡 Bringing you back to the main branch       "
$STAGEC = " ✅ COMPLETE!                                  "

# create the progress bar
$PROG0 = "[=>--------------------------------------] 0%"
$PROG1 = "[=>--------------------------------] 20%"
$PROG2 = "[=>------------------------] 40%"
$PROG3 = "[=>----------------] 60%"
$PROG4 = "[=>------------] 70%"
$PROG5 = "[=>--------] 80%"
$PROG6 = "[=>] 99%"
$COMPL = "[] 100%"

# print the first stage in purple
Write-Host " ${STAGE0}" -ForegroundColor Magenta
# print the progress bar
Write-Host " ${PROG0}" -ForegroundColor Magenta

# checkout the Beor branch if there is a problem then use the quit function
gitstep | Out-Null || quit "Could not checkout Beor branch"

# print the second stage in purple and delete the first stage
Write-Host "`r${STAGE1}" -ForegroundColor Magenta
# print the progress bar
Write-Host " ${PROG1}" -ForegroundColor Magenta
# do the version stuff
versionstep || quit "Failed to update VERSION"

# print the third stage in purple
Write-Host "`r${STAGE2}" -ForegroundColor Magenta
# print the progress bar
Write-Host " ${PROG2}" -ForegroundColor Magenta
# login to ECR and throw away output to avoid printing it
aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin 739810740537.dkr.ecr.us-east-1.amazonaws.com | Out-Null || quit "Failed to login to ECR"

# print the fourth stage in purple
Write-Host "`r${STAGE3}" -ForegroundColor Magenta
# print the progress bar
Write-Host " ${PROG3}" -ForegroundColor Magenta
# build the docker image and throw away output to avoid printing it
docker build -t tavern-profile-beor . | Out-Null || quit "Failed to build docker image"

# print the fifth stage in purple
Write-Host "`r${STAGE4}" -ForegroundColor Magenta
# print the progress bar
Write-Host " ${PROG4}" -ForegroundColor Magenta
# tag the docker image and throw away output to avoid printing it
docker tag tavern-profile-beor:latest 739810740537.dkr.ecr.us-east-1.amazonaws.com/beor:${MAJOR}.${MINOR} | Out-Null || quit "Failed to tag docker image"

# print the sixth stage in purple
Write-Host "`r${STAGE5}" -ForegroundColor Magenta
# print the progress bar
Write-Host " ${PROG5}" -ForegroundColor Magenta
# push the docker image and throw away output to avoid printing it
docker push 739810740537.dkr.ecr.us-east-1.amazonaws.com/beor:${MAJOR}.${MINOR} | Out-Null || quit "Failed to push docker image to ECR"

# print the seventh stage in purple
Write-Host "`r${STAGE6}" -ForegroundColor Magenta
# print the progress bar
Write-Host " ${PROG6}" -ForegroundColor Magenta
# checkout the main branch
returntoMain || quit "Failed to checkout main branch"

# print the complete stage in green
Write-Host "`r${STAGEC}" -ForegroundColor Green
# print the progress bar
Write-Host " ${COMPL}" -ForegroundColor Green

# print "Tavern Profile Pushed as Beor:${MAJOR}.${MINOR}"  in green
Write-Host " 🍻 Tavern Profile Pushed as Beor:${MAJOR}.${MINOR}" -ForegroundColor Green





