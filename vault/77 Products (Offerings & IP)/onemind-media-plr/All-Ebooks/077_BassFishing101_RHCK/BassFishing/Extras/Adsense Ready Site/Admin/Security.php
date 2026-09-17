<?php include("Session.php"); ?>
<?php

	$strSecurity					= "disabled";
	
	if (isset($_POST['UserName']))
		{						
			$_SESSION['UserName']	= $_POST['UserName'];
			$_SESSION['Password']	= $_POST['Password'];
		
			header('location: Write.php');
		}
?>
	<?php include("Header.php"); ?>
		<form method="post" action="Security.php" name="frmData">
			<center>
				<table class="tblForm">
					<tr>
						<td align="center" style="background: url(http://www.niche-mania.com/Images/NM_006.jpg); color: #FFFFFF; font-weight: bold; height: 24px; text-align: center; vertical-align: middle;" colspan="2">
							Admin Security 
						</td>				
					<tr>
					<tr>
						<td colspan="2">
							&nbsp;
						</td>				
					<tr>
					<tr>
						<td class="tdFieldName">
							User Name
						</td>				
						<td class="tdFieldValue">
							<input class="iText" type="text" style="width: 100%;" name="UserName" value="<?php echo $_SESSION['UserName']; ?>">
						</td>				
					<tr>
					<tr>
						<td class="tdFieldName">
							Password
						</td>				
						<td class="tdFieldValue">
							<input class="iText" type="text" style="width: 100%;" name="Password" value="<?php echo $_SESSION['Password']; ?>">
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
							<input class="iButton" onClick="this.disabled='true'; frmData.submit();" type="submit" value=" Save Security Settings ">					
						</td>				
					<tr>
				</table>
			</center>
		</form>
	<?php include("Footer.php"); ?>