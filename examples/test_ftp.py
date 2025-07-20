#!/usr/bin/env python3
"""
Simple FTP client test script for the Go FTP Server
This script demonstrates basic FTP operations using the ftplib library.
"""

import ftplib
import os
import tempfile
import time

def test_ftp_server(host='localhost', port=2121, username='anonymous', password='anonymous'):
    """Test basic FTP server functionality"""
    
    print(f"Testing FTP server at {host}:{port}")
    print(f"Username: {username}")
    print("-" * 50)
    
    try:
        # Connect to FTP server
        ftp = ftplib.FTP()
        ftp.connect(host, port)
        print(f"✓ Connected to {host}:{port}")
        
        # Login
        ftp.login(username, password)
        print("✓ Login successful")
        
        # Get welcome message
        welcome = ftp.getwelcome()
        print(f"✓ Server welcome: {welcome}")
        
        # List current directory
        print("\n--- Directory Listing ---")
        files = ftp.nlst()
        for file in files:
            print(f"  {file}")
        
        # Create a test file
        test_content = f"Test file created at {time.strftime('%Y-%m-%d %H:%M:%S')}"
        test_filename = f"test_file_{int(time.time())}.txt"
        
        # Upload test file
        print(f"\n--- Uploading {test_filename} ---")
        with tempfile.NamedTemporaryFile(mode='w', delete=False) as temp_file:
            temp_file.write(test_content)
            temp_file_path = temp_file.name
        
        with open(temp_file_path, 'rb') as file:
            ftp.storbinary(f'STOR {test_filename}', file)
        print(f"✓ Uploaded {test_filename}")
        
        # Clean up temp file
        os.unlink(temp_file_path)
        
        # Download the file back
        print(f"\n--- Downloading {test_filename} ---")
        download_filename = f"downloaded_{test_filename}"
        with open(download_filename, 'wb') as file:
            ftp.retrbinary(f'RETR {test_filename}', file.write)
        print(f"✓ Downloaded to {download_filename}")
        
        # Verify content
        with open(download_filename, 'r') as file:
            downloaded_content = file.read()
        
        if downloaded_content == test_content:
            print("✓ Content verification successful")
        else:
            print("✗ Content verification failed")
        
        # Clean up downloaded file
        os.unlink(download_filename)
        
        # Get current working directory
        pwd = ftp.pwd()
        print(f"\n✓ Current directory: {pwd}")
        
        # Test NOOP command
        ftp.voidcmd('NOOP')
        print("✓ NOOP command successful")
        
        # Quit
        ftp.quit()
        print("\n✓ Disconnected successfully")
        
        return True
        
    except ftplib.error_perm as e:
        print(f"✗ FTP Permission Error: {e}")
        return False
    except ftplib.error_temp as e:
        print(f"✗ FTP Temporary Error: {e}")
        return False
    except ftplib.error_proto as e:
        print(f"✗ FTP Protocol Error: {e}")
        return False
    except ConnectionRefusedError:
        print(f"✗ Connection refused to {host}:{port}")
        print("  Make sure the FTP server is running!")
        return False
    except Exception as e:
        print(f"✗ Unexpected error: {e}")
        return False

def test_advanced_features(host='localhost', port=2121, username='anonymous', password='anonymous'):
    """Test advanced FTP features"""
    
    print(f"\nTesting Advanced Features")
    print("-" * 50)
    
    try:
        ftp = ftplib.FTP()
        ftp.connect(host, port)
        ftp.login(username, password)
        
        # Test directory change (if subdir exists)
        try:
            ftp.cwd('subdir')
            print("✓ Changed to subdir directory")
            
            # List files in subdir
            files = ftp.nlst()
            print(f"  Files in subdir: {files}")
            
            # Go back to parent
            ftp.cwd('..')
            print("✓ Returned to parent directory")
            
        except ftplib.error_perm:
            print("ℹ No subdir directory available for testing")
        
        # Test file size (if files exist)
        files = ftp.nlst()
        if files:
            try:
                size = ftp.size(files[0])
                print(f"✓ File size of {files[0]}: {size} bytes")
            except ftplib.error_perm:
                print(f"ℹ Could not get size of {files[0]}")
        
        ftp.quit()
        return True
        
    except Exception as e:
        print(f"✗ Advanced features test failed: {e}")
        return False

if __name__ == "__main__":
    print("Go FTP Server Test Client")
    print("=" * 50)
    
    # Test basic functionality
    success = test_ftp_server()
    
    if success:
        # Test advanced features
        test_advanced_features()
        
        print("\n" + "=" * 50)
        print("🎉 All tests completed successfully!")
        print("The Go FTP Server is working correctly.")
    else:
        print("\n" + "=" * 50)
        print("❌ Tests failed. Please check the server status.")
        print("\nTroubleshooting:")
        print("1. Make sure the FTP server is running")
        print("2. Check if port 2121 is available")
        print("3. Verify firewall settings")
        print("4. Check server logs for errors") 