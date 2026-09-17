<?php include("../Includes/Start.php"); ?>
<html>
	<head>
		<title><?=$Category;?> Legal Information</title>
		<meta http-equiv="Content-Type"	content="text/html; charset=windows-1252">
		<meta name="Creator"			content="Niche Mania - http://www.niche-maniacs.com">
		<meta name="Copyright"			content="<?=$CopyRight;?>">
		<meta name="Description"		content="<?=$Category;?> Legal Information">
		<meta name="Keywords"			content="<?=$Keywords;?>">
		<meta name="Distribution"		content="global">
		<meta name="Publisher"			content="<?=$Domain;?>">
		<meta name="Rating"				content="General">
		<meta name="Revisit-after"		content="5 days">
		<meta name="Robots"				content="index,follow">
		<link href="../Includes/Styles.css" rel="stylesheet" type="text/css">
		<script src="../Includes/JavaScript.js"></script>		
	</head>
	<body>
		<table align="center" cellpadding="0" cellspacing="0" class="tblMain">
			<tr>
				<td class="tdHeader" colspan="2">
					<h1>
						<?=$MainTitle;?>
					</h1>
					<h3>
						<?=$SubTitle;?>
					</h3>
				</td>
			</tr>
			<tr>
				<td class="tdRow" colspan="2" >
					<?php include("../Includes/Menu.php"); ?>
				</td>
			</tr>
			<tr>
				<td class="tdContent" style="border-right-width: 0; width: 780;">
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
					<?php include("../Includes/Legal.php"); ?>
				</td>
			</tr>
			<tr>
				<td class="tdRow" colspan="2">
					&copy; <?=date("Y");?>, <a href="<?=$Domain;?>"><?=$SiteName;?></a> - All Rights Reserved Worldwide | <b><?=$Category;?> Legal Information</b>
				</td>
			</tr>
		</table>
		<?php include("../Includes/Footer.php"); ?>
		<?php include("../Includes/AdTracker.php"); ?>
	</body>
</html>

