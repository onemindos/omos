<?php
	
	include('../Includes/Configure.php');
		
	if (isset($_POST['UserName']) and isset($_POST['Password']))
		{						
			if ($_POST['UserName'] != '' and $_POST['Password'] != '')
				{
					if ($_POST['UserName'] == $UserName and $_POST['Password'] == $Password)
						{
							session_start();
	
							$_SESSION['UserName']		= $UserName;
							$_SESSION['Password']		= $Password;
							$_SESSION['Domain']			= $Domain;
							$_SESSION['SiteName']		= stripslashes($SiteName);
							$_SESSION['Category']		= $Category;
							$_SESSION['NewsFeed']		= $NewsFeed;
							$_SESSION['MainTitle']		= stripslashes($MainTitle);
							$_SESSION['SubTitle']		= stripslashes($SubTitle);
							$_SESSION['Keywords']		= $Keywords;							
							$_SESSION['GooglePubID']	= $GooglePubID;
							$_SESSION['GoogleChannel']	= $GoogleChannel;
							$_SESSION['GoogleLinkColor']= $GoogleLinkColor;
							$_SESSION['GoogleURLColor']	= $GoogleURLColor;														
							$_SESSION['AmazonID']		= $AmazonID;
							$_SESSION['AmazonMode']		= $AmazonMode;
							$_SESSION['AmazonSearch']	= $AmazonSearch;
							$_SESSION['AmazonLinkColor']= $AmazonLinkColor;
							$_SESSION['AmazonTextColor']= $AmazonTextColor;
							$_SESSION['DisplayAmazon']	= $DisplayAmazon;							
							$_SESSION['CopyRight']		= $CopyRight;
							$_SESSION['Name']			= stripslashes($Name);
							$_SESSION['Address1']		= stripslashes($Address1);
							$_SESSION['Address2']		= stripslashes($Address2);
							$_SESSION['City']			= $City;
							$_SESSION['State']			= $State;
							$_SESSION['ZipCode']		= $ZipCode;
							$_SESSION['Country']		= $Country;
							$_SESSION['Telephone']		= $Telephone;
							$_SESSION['EmailAddress']	= $EmailAddress;
							$_SESSION['ShowGoogle']		= $ShowGoogle;
							$_SESSION['ShowFireFox']	= $ShowFireFox;														
							$_SESSION['ShowSearch']		= $ShowSearch;														
							
							header('location: Site.php');
					}
				else
					{
						header('location: index.php?Status=E');
					}
				}
			else
				{
					header('location: index.php?Status=E');
				}
		}		
	else
		{
			header('location: index.php?Status=E');
		}

 ?>