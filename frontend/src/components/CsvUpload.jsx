import { useState } from "react";
import axios from "axios";
import PropTypes from "prop-types";

function CsvUpload({ token, endpoint, onSuccess, onError,uploadtype }) {
  const [csvFile, setCsvFile] = useState(null);

  const handleCsvUpload = (e) => {
    setCsvFile(e.target.files[0]);
  };

  const uploadCsv = async () => {
    if (!csvFile) {
      onError("Please select a CSV file.");
      return;
    }

    const formData = new FormData();
    formData.append("csvFile", csvFile);

    try {
      const response = await axios.post(endpoint, formData, {
        headers: {
          Authorization: `Bearer ${token}`,
          "Content-Type": "multipart/form-data",
        },
      });

      if (response.status === 200) {
        onSuccess(`CSV file uploaded successfully to ${endpoint}.`);
      }
    } catch (error) {
      onError(`Failed to upload CSV file to ${endpoint}: ${error.message}`);
    }
  };

  return (
    <div className="p-2 m-0">
      <input type="file" accept=".csv" onChange={handleCsvUpload} />
      <button
        onClick={uploadCsv}
        className="bg-blue-500 hover:bg-blue-600 text-white rounded p-2"
      >
        Upload {uploadtype} CSV
      </button>
    </div>
  );
}
CsvUpload.propTypes = {
  token: PropTypes.string.isRequired,
  endpoint: PropTypes.string.isRequired,
  onSuccess: PropTypes.func.isRequired,
  onError: PropTypes.func.isRequired,
  uploadtype: PropTypes.string.isRequired,
};

export default CsvUpload;
