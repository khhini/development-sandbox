import os
from typing import Dict, List, Union
import gspread
from google.oauth2.service_account import Credentials

SCOPES = [
    "https://www.googleapis.com/auth/spreadsheets",
    "https://www.googleapis.com/auth/drive",
]

SHEETS_URL = "https://docs.google.com/spreadsheets/d/1PwMDOtQ-hw877UgeMsT1MiuhVQXdYGWuNHp49Us5NzY/edit?gid=218693089#gid=218693089"


def get_tracker_data(
    client: gspread.Client, sheet_url
) -> List[Dict[str, Union[int, float, str]]]:
    sheet = client.open_by_url(sheet_url)
    worksheet = sheet.worksheet("Tracker")
    return worksheet.get_all_records()


def main():
    service_account_key = os.getenv("GOOGLE_APPLICATION_CREDENTIALS")

    credentials = Credentials.from_service_account_file(
        service_account_key, scopes=SCOPES
    )

    client = gspread.authorize(credentials)

    tracker_data = get_tracker_data(client, SHEETS_URL)

    print(tracker_data[0])


if __name__ == "__main__":
    main()
