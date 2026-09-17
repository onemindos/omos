<?php include("Session.php"); ?>
<?php

	if (isset($_POST['FullName']))
		{						

			$msgText .= "R E Q U E S T   F O R M   D A T A\n\n";
			
			$msgText  = "      Customer Name: ".$_POST['FullName']."\n";
			$msgText .= "      Email Address: ".$_POST['EmailAddress']."\n";
			$msgText .= "        Description: ".$_POST['Description']."\n\n";			

			$msgText .= "S E R V E R   V A R I A B L E S\n\n";
			
			$msgText .= "Customer IP Address: ".$_SERVER["REMOTE_ADDR"]."\n";
			$msgText .= "   Installed Domain: ".$_SERVER["HTTP_HOST"]."\n";
			$msgText .= "      Document Root: ".$_SERVER["DOCUMENT_ROOT"]."\n\n";

			$msgText .= "S E S S I O N   V A R I A B L E S\n\n";			

			$msgText .= "Site Configuration...\n\n";
			
			$msgText .= "             Domain: ".$_SESSION['Domain']."\n";
			$msgText .= "          Site Name: ".$_SESSION['SiteName']."\n";
			$msgText .= "         Main Title: ".$_SESSION['MainTitle']."\n";
			$msgText .= "          Sub Title: ".$_SESSION['SubTitle']."\n\n";
			
			$msgText .= "Adsense Settings...\n\n";

			$msgText .= "     G Publisher ID: ".$_SESSION['GooglePubID']."\n";
			$msgText .= "     G Channel Code: ".$_SESSION['GoogleChannel']."\n";
			$msgText .= "       G Link Color: ".$_SESSION['GoogleLinkColor']."\n";
			$msgText .= "        G URL Color: ".$_SESSION['GoogleURLColor']."\n\n";
			
			$msgText .= "Contact Information...\n\n";
			
			$msgText .= "              Name: ".$_SESSION['Name']."\n";
			$msgText .= "         Address 1: ".$_SESSION['Address1']."\n";
			$msgText .= "         Address 2: ".$_SESSION['Address2']."\n";
			$msgText .= "              City: ".$_SESSION['City']."\n";
			$msgText .= "             State: ".$_SESSION['State']."\n";
			$msgText .= "          Zip Code: ".$_SESSION['ZipCode']."\n";
			$msgText .= "           Country: ".$_SESSION['Country']."\n";
			$msgText .= "         Telephone: ".$_SESSION['Telephone']."\n";
			$msgText .= "     Email Address: ".$_SESSION['EmailAddress']."\n\n";

			$msgText .= "Login Information...\n\n";

			$msgText .= "        User Name: ".$_SESSION['UserName']."\n";
			$msgText .= "         Password: ".$_SESSION['Password']."\n";			
			$msgText .= "        Login URL: http://".$_SERVER["HTTP_HOST"]."/Admin\n";
						
			mail("mike@anamik.com", "Niche Mania Help Desk Support Request", $msgText, "From: webmaster@".$_SERVER['HTTP_HOST']);
		}
?>
	<?php include("Header.php"); ?>
		<form method="post" action="Support.php" name="frmData">
			<center>
				<table class="tblForm">
					<tr>
						<td align="center" style="background: url(http://www.niche-mania.com/Images/NM_006.jpg); color: #FFFFFF; font-weight: bold; height: 24px; text-align: center; vertical-align: middle;" colspan="2">
							HELP DESK SUPPORT REQUEST &nbsp; 
						</td>				
					<tr>
					<tr>
						<td colspan="2">
							<small>
								&nbsp;
							</small>
						</td>				
					<tr>
					<tr>
						<td class="tdFieldName">
							Full Name
						</td>				
						<td class="tdFieldValue">
							<input class="iText" type="text" size="50" name="FullName" style="width: 100%;" value="<?php echo $_SESSION['Name']; ?>">
						</td>				
					<tr>
					<tr>
						<td class="tdFieldName">
							Email Address
						</td>				
						<td class="tdFieldValue">
							<input class="iText" type="text" size="50" name="EmailAddress" style="width: 100%;" value="<?php echo $_SESSION['EmailAddress']; ?>">
						</td>				
					<tr>
					<tr>
						<td class="tdFieldName">
							Problem Description
						</td>				
						<td class="tdFieldValue">
							<textarea class="taText" name="Description" rows="10" style="width: 100%;"></textarea>
						</td>				
					<tr>
					<tr>
						<td colspan="2">
							<hr size="0">
						</td>				
					<tr>
					<tr>
						<td class="tdFieldName">
						</td>				
						<td class="tdFieldValue">
							<input class="iButton" onClick="this.disabled='true'; frmData.submit();" type="submit" value=" Submit To Help Desk ">					
						</td>				
					<tr>
				</table>
			</center>
		</form>
	<?php include("Footer.php"); ?>