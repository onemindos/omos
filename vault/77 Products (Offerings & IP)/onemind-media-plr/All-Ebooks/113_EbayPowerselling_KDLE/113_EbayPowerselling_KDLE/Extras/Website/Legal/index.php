<?php include("../Includes/Configure.php"); ?>
<html>
	<head>
		<title><?=$Category;?> Legal Information</title>
		<META http-equiv="Content-Type" content="text/html; charset=windows-1252">
		<meta name="Creator"			content="Kotzakoliou, SSA - http://www.niche-mania.com">
		<meta name="Copyright"			content="<?=$CopyRight;?>">
		<meta name="Description"		content="<?=$Category;?> Legal Information">
		<meta name="Keywords"			content="<?=$Keywords;?>">
		<meta name="Distribution"		content="global">
		<meta name="Publisher"			content="<?=$Domain;?>">
		<meta name="Rating"				content="General">
		<meta name="Revisit-after"		content="5 days">
		<meta name="Robots"				content="index,follow">
		<link href="../Includes/Styles.css" rel="stylesheet" type="text/css">
	</head>
	<body>
		<table align="center" cellpadding="0" cellspacing="0">
			<tr>
				<td class="tdHeader" colspan="3">
					<p style="margin-bottom: 5; margin-left: 10;">
						<?php include("../Includes/Menu.php"); ?>
					</p>
				</td>
			</tr>
			<tr>
				<td colspan="3">
					<?php include("../Includes/Google_728x15_TextLinks.php"); ?>
					<hr>
				</td>
			</tr>
			<tr>
				<td class="tdleft">
					<br>
					<?php include("../Includes/Navigation.php"); ?>
				</td>
				<td class="tdContent">
					<br>
					<p>
						<b>
							<?=$SiteName;?> Legal Information
						</b>
						<br>
						<?php
							if ($Name != "")
								{
									echo "<br>&nbsp;&nbsp;",stripslashes($Name);
								}
							if ($Address1 != "")
								{
									echo "<br>&nbsp;&nbsp;".stripslashes($Address1);
								}
							if ($Address2 != "")
								{
									echo "<br>&nbsp;&nbsp;".stripslashes($Address2);
								}
							if ($City != "")
								{
									echo "<br>&nbsp;&nbsp;".stripslashes($City);
								}
							if ($State != "")
								{
									echo ", $State";
								}
							if ($ZipCode != "")
								{
									echo ", $ZipCode";
								}
							if ($Country != "")
								{
									echo "<br>&nbsp;&nbsp;$Country";
								}
							if ($Telephone != "")
								{
									echo "<br>&nbsp;&nbsp;$Telephone";
								}
							if ($EmailAddress != "")
								{
									echo "<br>&nbsp;&nbsp;<a href=mailto:$EmailAddress>$EmailAddress</a>";
								}
						?>
					</p>
					<?php include("../Includes/Google_336x280.php"); ?>
					<?php include("../Includes/Legal.php"); ?>
				</td>
				<td class="tdRight">
					<?php include("../Includes/Google_160x600.php"); ?>					
				</td>
			</tr>
			<tr>
				<td class="tdFooter" colspan="3">
					&copy; 2005, <a href="<?=$Domain;?>"><?=$SiteName;?></a> - All Rights Reserved Worldwide | <a href="../Legal/index.php"><?=$Category;?> Legal Information</a>
				</td>
			</tr>
		</table>
		<?php include("../Includes/AdTracker.php"); ?>
	</body>
</html>
