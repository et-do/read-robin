import React, { useState, useEffect } from "react";
import { db } from "./firebase";
import { doc, getDoc, deleteDoc } from "firebase/firestore";
import "./AttemptPage.css";

function AttemptPage({ user, activePersona, contentID, attemptID, setPage }) {
  const [attempt, setAttempt] = useState(null);
  const [loading, setLoading] = useState(true);
  const [quizTitle, setQuizTitle] = useState("");
  const [contentType, setContentType] = useState("");
  const [error, setError] = useState(null);
  const [showConfirmation, setShowConfirmation] = useState(false);

  useEffect(() => {
    const fetchAttempt = async () => {
      if (user && activePersona && contentID && attemptID) {
        const quizRef = doc(
          db,
          "users",
          user.uid,
          "personas",
          activePersona.id,
          "quizzes",
          contentID
        );
        const quizDoc = await getDoc(quizRef);
        if (quizDoc.exists()) {
          const data = quizDoc.data();
          setQuizTitle(data.title);
          setContentType(data.content_type);
        }

        const attemptRef = doc(
          db,
          "users",
          user.uid,
          "personas",
          activePersona.id,
          "quizzes",
          contentID,
          "attempts",
          attemptID
        );
        const attemptDoc = await getDoc(attemptRef);
        if (attemptDoc.exists()) {
          setAttempt(attemptDoc.data());
        } else {
          console.error("No such document!");
        }
        setLoading(false);
      }
    };
    fetchAttempt();
  }, [user, activePersona, contentID, attemptID]);

  const handleDeleteClick = () => {
    setShowConfirmation(true);
  };

  const handleConfirmDelete = async () => {
    setLoading(true);
    setShowConfirmation(false);
    try {
      const attemptDocRef = doc(
        db,
        "users",
        user.uid,
        "personas",
        activePersona.id,
        "quizzes",
        contentID,
        "attempts",
        attemptID
      );
      await deleteDoc(attemptDocRef);
      setPage("selection");
    } catch (error) {
      console.error("Error deleting attempt:", error);
      setError("Error deleting attempt: " + error.message);
    } finally {
      setLoading(false);
    }
  };

  const handleCancelDelete = () => {
    setShowConfirmation(false);
  };

  const getScoreClass = (score) => {
    if (score <= 50) return "red";
    if (score >= 80) return "green";
    return "";
  };

  const getStatusClass = (status) => {
    return status === "Correct" ? "green" : "red";
  };

  if (loading) {
    return <div>Loading...</div>;
  }

  if (!attempt) {
    return <div>No data found</div>;
  }

  return (
    <div className="attempt-page">
      <button
        className="back-button"
        onClick={() => setPage("performanceHistory")}
      >
        Back
      </button>
      <button className="delete-attempt-button" onClick={handleDeleteClick}>
        &times;
      </button>
      <h2>{quizTitle}</h2>
      <p className="content-type">Content Type: {contentType}</p>
      <p className="score-text">
        Score:{" "}
        <span className={`score-value ${getScoreClass(attempt.score)}`}>
          {attempt.score}%
        </span>
      </p>
      <ul>
        {attempt.responses.map((response, index) => (
          <li key={index}>
            <div className="response-item">
              <span className="response-title">Question:</span>
              <span>{response.question}</span>
            </div>
            <div className="response-item">
              <span className="response-title">Correct Answer:</span>
              <span>{response.answer}</span>
            </div>
            <div className="response-item">
              <span className="response-title">Your Response:</span>
              <span>{response.userResponse}</span>
            </div>
            <div className="response-item">
              <span className="response-title">Status:</span>
              <span
                className={`status-value ${getStatusClass(response.status)}`}
              >
                {response.status}
              </span>
            </div>
            <div className="response-item">
              <span className="response-title">Reference:</span>
              <span>{response.reference}</span>
            </div>
            {index < attempt.responses.length - 1 && <hr />}
          </li>
        ))}
      </ul>
      <p className="attempt-id">Attempt ID: {attempt.attemptID}</p>
      {error && <div style={{ color: "red" }}>{error}</div>}
      {showConfirmation && (
        <div className="confirmation-modal">
          <div className="confirmation-content">
            <h3>Are you sure you want to permanently delete this attempt?</h3>
            <p>This action cannot be undone.</p>
            <button
              className="confirm-button"
              onClick={handleConfirmDelete}
              disabled={loading}
            >
              {loading ? "Deleting..." : "Yes, Delete"}
            </button>
            <button className="cancel-button" onClick={handleCancelDelete}>
              Cancel
            </button>
            {error && <p style={{ color: "red" }}>{error}</p>}
          </div>
        </div>
      )}
    </div>
  );
}

export default AttemptPage;
