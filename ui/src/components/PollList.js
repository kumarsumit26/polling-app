import React, { useEffect, useState } from 'react';
import axios from 'axios';
import Modal from './Modal';
import './PollList.css';

const PollList = () => {
  const [polls, setPolls] = useState([]);
  const [newPoll, setNewPoll] = useState('');
  const [errorMessage, setErrorMessage] = useState('');
  const [showModal, setShowModal] = useState(false);

  useEffect(() => {
    const fetchPolls = async () => {
      try {
        const token = localStorage.getItem('token');
        const response = await axios.get('/polls', {
          headers: { Authorization: `Bearer ${token}` },
        });
        setPolls(response.data.sort((a, b) => a.id - b.id));
      } catch (error) {
        console.error('Error fetching polls:', error);
      }
    };

    fetchPolls();
  }, []);

  const handleAddPoll = async (e) => {
    e.preventDefault();
    try {
      const token = localStorage.getItem('token');
      await axios.post(
        '/polls',
        { question: newPoll },
        {
          headers: { Authorization: `Bearer ${token}`, 'Content-Type': 'application/json' },
        }
      );
      setNewPoll('');
      // Refresh the poll list
      const response = await axios.get('/polls', {
        headers: { Authorization: `Bearer ${token}` },
      });
      setPolls(response.data.sort((a, b) => a.id - b.id));
    } catch (error) {
      console.error('Error adding poll:', error);
    }
  };

  const handleVote = async (pollId, vote) => {
    try {
      const token = localStorage.getItem('token');
      await axios.post(
        `/polls/${pollId}/vote`,
        { vote },
        {
          headers: { Authorization: `Bearer ${token}`, 'Content-Type': 'application/json' },
        }
      );
      // Refresh the poll list
      const response = await axios.get('/polls', {
        headers: { Authorization: `Bearer ${token}` },
      });
      setPolls(response.data.sort((a, b) => a.id - b.id));
      setErrorMessage('');
    } catch (error) {
      if (error.response && error.response.status === 409) {
        setErrorMessage('User has already voted');
        setShowModal(true);
      } else {
        console.error(`Error voting on poll ${pollId}:`, error);
      }
    }
  };

  const handleCloseModal = () => {
    setShowModal(false);
    setErrorMessage('');
  };

  return (
    <div className="container">
      <h2>Poll List</h2>
      <form onSubmit={handleAddPoll} className="add-poll-form">
        <div className="form-group">
          <label>Create a new poll below:</label>
          <input
            type="text"
            value={newPoll}
            onChange={(e) => setNewPoll(e.target.value)}
            required
          />
        </div>
        <button type="submit" className="btn">Add Poll</button>
      </form>
      <ul className="poll-list">
        {polls.map((poll) => (
          <li key={poll.id} className="poll-item">
            <p className="poll-question">{poll.question}</p>
            <p className="poll-counts">Yes: {poll.yes_count} | No: {poll.no_count}</p>
            <div className="poll-buttons">
              <button onClick={() => handleVote(poll.id, 'yes')} className="btn">Yes</button>
              <button onClick={() => handleVote(poll.id, 'no')} className="btn">No</button>
            </div>
          </li>
        ))}
      </ul>
      <Modal show={showModal} handleClose={handleCloseModal}>
        <p>{errorMessage}</p>
      </Modal>
    </div>
  );
};

export default PollList;