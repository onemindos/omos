<?php include("Session.php"); ?>
<?php

	if (isset($_POST['SiteName']))
		{						
			$_SESSION['SiteName']	= $_POST['SiteName'];
			$_SESSION['MainTitle']	= $_POST['MainTitle'];
			$_SESSION['SubTitle']	= $_POST['SubTitle'];
		
			header('location: Write.php');
		}
		
	if (strstr(file_get_contents("../Includes/Template.htm"), "MainTitle;"))
		{
			$strMainTitle = true;
		}
	if (strstr(file_get_contents("../Includes/Template.htm"), "SubTitle;"))
		{
			$strSubTitle = true;
		}
?>
	<?php include("Header.php"); ?>
		<form method="post" action="Site.php" name="frmData">
			<center>
				<table class="tblForm">
					<tr>
						<td align="center" style="background: url(http://www.niche-mania.com/Images/NM_006.jpg); color: #FFFFFF; font-weight: bold; height: 24px; text-align: center; vertical-align: middle;" colspan="2">
							SITE CONFIGURATION&nbsp; 
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
							Domain Name
						</td>				
						<td class="tdFieldValue">
							<input disabled class="iText" type="text" size="50" name="Domain" value="<?php echo $_SESSION['Domain']; ?>">
						</td>				
					<tr>
					<tr>
						<td class="tdFieldName">
							Site Name
						</td>				
						<td class="tdFieldValue">
							<input class="iText" type="text" size="50" name="SiteName" value="<?php echo $_SESSION['SiteName']; ?>">
						</td>				
					<tr>
					<?php 
						if ($strMainTitle)
							{
					?>
								<tr>
									<td class="tdFieldName">
										Main Title
									</td>				
									<td class="tdFieldValue">
										<input class="iText" type="text" size="50" name="MainTitle" value="<?php echo $_SESSION['MainTitle']; ?>">
									</td>				
								<tr>
					<?php
							}
						if ($strSubTitle)
							{
					?>
								<tr>
									<td class="tdFieldName">
										Sub Title
									</td>				
									<td class="tdFieldValue">
										<input class="iText" type="text" size="50" name="SubTitle" value="<?php echo $_SESSION['SubTitle']; ?>">
									</td>				
								<tr>
					<?php
							}
					?>
					<tr>
						<td colspan="2">
							<hr size="0">
						</td>				
					<tr>
					<tr>
						<td class="tdFieldName">
						</td>				
						<td class="tdFieldValue">
							<input class="iButton" onClick="this.disabled='true'; frmData.submit();" type="submit" value=" Save Site Configuration ">
						</td>				
					<tr>
				</table>
			</center>
		</form>
	<?php include("Footer.php"); ?>