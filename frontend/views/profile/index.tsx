"use client";

import { useState, useEffect, useCallback, useRef } from "react";
import { Camera, Trash2, Plus, X } from "lucide-react";
import { Card, Button, Input } from "@/shared/ui";
import { TabBar } from "@/shared/ui";
import { RouteGuard } from "@/features/auth/ui/RouteGuard";
import { getProfile, updateProfile, type Profile } from "@/entities/user/api/profileRepository";
import {
  getMyPhotos, uploadPhoto, deletePhoto, uploadAvatar, type Photo,
} from "@/entities/photo/api/photoRepository";
import styles from "./profile.module.css";

interface Goals { calorie: string; protein: string; carbs: string; fat: string; }

function toGoals(p: Profile): Goals {
  return {
    calorie: String(p.calorie_goal || ""),
    protein: String(p.protein_goal || ""),
    carbs:   String(p.carbs_goal   || ""),
    fat:     String(p.fat_goal     || ""),
  };
}

function ProfileScreen() {
  const [profile, setProfile] = useState<Profile | null>(null);
  const [photos, setPhotos]   = useState<Photo[]>([]);
  const [goals, setGoals]     = useState<Goals>({ calorie: "", protein: "", carbs: "", fat: "" });
  const [editingGoals, setEditingGoals] = useState(false);
  const [saving, setSaving]   = useState(false);
  const [savedMsg, setSavedMsg] = useState(false);
  const [err, setErr]         = useState<string | null>(null);

  const photoInput  = useRef<HTMLInputElement>(null);
  const avatarInput = useRef<HTMLInputElement>(null);

  const load = useCallback(async () => {
    setErr(null);
    try {
      const [p, ph] = await Promise.all([getProfile(), getMyPhotos()]);
      setProfile(p);
      setGoals(toGoals(p));
      setPhotos(ph);
    } catch (e) {
      setErr((e as Error).message);
    }
  }, []);

  useEffect(() => { load(); }, [load]);

  const handleSaveGoals = async () => {
    if (!profile) return;
    setSaving(true);
    setErr(null);
    try {
      // Echo back weight/direction so the upsert doesn't wipe them.
      const updated = await updateProfile({
        calorie_goal:   Number(goals.calorie) || 0,
        protein_goal:   Number(goals.protein) || 0,
        carbs_goal:     Number(goals.carbs)   || 0,
        fat_goal:       Number(goals.fat)     || 0,
        current_weight: profile.current_weight,
        goal_weight:    profile.goal_weight,
        direction:      profile.direction,
      });
      setProfile(updated);
      setSavedMsg(true);
      setTimeout(() => { setSavedMsg(false); setEditingGoals(false); }, 1200);
    } catch (e) {
      setErr((e as Error).message);
    } finally {
      setSaving(false);
    }
  };

  const handlePhotoSelect = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    e.target.value = ""; // allow re-selecting the same file
    if (!file) return;
    try {
      await uploadPhoto(file);
      await load();
    } catch (e) {
      setErr((e as Error).message);
    }
  };

  const handleAvatarSelect = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    e.target.value = "";
    if (!file) return;
    try {
      const { avatar_url } = await uploadAvatar(file);
      setProfile((p) => (p ? { ...p, avatar_url } : p));
    } catch (e) {
      setErr((e as Error).message);
    }
  };

  const handleDeletePhoto = async (id: number) => {
    try {
      await deletePhoto(id);
      setPhotos((ps) => ps.filter((p) => p.id !== id));
    } catch (e) {
      setErr((e as Error).message);
    }
  };

  const name = profile?.first_name || profile?.username || "Your profile";

  return (
    <>
      <main className={styles.root}>
        <h1 className={styles.title}>Profile</h1>

        {err && <p className={styles.errMsg}>{err}</p>}

        {/* Identity */}
        <Card>
          <div className={styles.identity}>
            <button
              className={styles.avatarWrap}
              onClick={() => avatarInput.current?.click()}
              aria-label="Change avatar"
            >
              {profile?.avatar_url ? (
                // eslint-disable-next-line @next/next/no-img-element
                <img className={styles.avatar} src={profile.avatar_url} alt="" />
              ) : (
                <div className={styles.avatarFallback}>{name.charAt(0).toUpperCase()}</div>
              )}
              <span className={styles.avatarBadge}><Camera size={14} /></span>
            </button>
            <div className={styles.identityText}>
              <p className={styles.name}>{name}</p>
              {profile?.username && <p className={styles.handle}>@{profile.username}</p>}
            </div>
          </div>
        </Card>

        {/* Goals */}
        <Card>
          <div className={styles.goalsHeader}>
            <p className={styles.sectionLabel}>Daily goals</p>
            {editingGoals ? (
              <button
                className={styles.iconLink}
                onClick={() => { setEditingGoals(false); if (profile) setGoals(toGoals(profile)); }}
                aria-label="Close editor"
              >
                <X size={18} />
              </button>
            ) : (
              <button className={styles.editLink} onClick={() => setEditingGoals(true)}>
                Edit
              </button>
            )}
          </div>

          {editingGoals ? (
            <>
              <div className={styles.goalsGrid}>
                <Input label="Calories" type="number" inputMode="numeric" value={goals.calorie}
                  onChange={(e) => setGoals({ ...goals, calorie: e.target.value })} rightSlot="kcal" />
                <Input label="Protein" type="number" inputMode="numeric" value={goals.protein}
                  onChange={(e) => setGoals({ ...goals, protein: e.target.value })} rightSlot="g" />
                <Input label="Carbs" type="number" inputMode="numeric" value={goals.carbs}
                  onChange={(e) => setGoals({ ...goals, carbs: e.target.value })} rightSlot="g" />
                <Input label="Fat" type="number" inputMode="numeric" value={goals.fat}
                  onChange={(e) => setGoals({ ...goals, fat: e.target.value })} rightSlot="g" />
              </div>
              <Button fullWidth onClick={handleSaveGoals} loading={saving} className={styles.saveBtn}>
                {savedMsg ? "Saved ✓" : "Save goals"}
              </Button>
            </>
          ) : (
            <div className={styles.goalsSummary}>
              <div className={styles.summaryItem}><span>Calories</span><strong>{goals.calorie || "—"} kcal</strong></div>
              <div className={styles.summaryItem}><span>Protein</span><strong>{goals.protein || "—"} g</strong></div>
              <div className={styles.summaryItem}><span>Carbs</span><strong>{goals.carbs || "—"} g</strong></div>
              <div className={styles.summaryItem}><span>Fat</span><strong>{goals.fat || "—"} g</strong></div>
            </div>
          )}
        </Card>

        {/* My photos */}
        <Card>
          <div className={styles.photosHeader}>
            <p className={styles.sectionLabel}>My photos</p>
            <button className={styles.addPhotoLink} onClick={() => photoInput.current?.click()}>
              <Plus size={16} /> Add
            </button>
          </div>
          {photos.length === 0 ? (
            <p className={styles.empty}>No photos yet. Tap + to post one.</p>
          ) : (
            <div className={styles.grid}>
              {photos.map((p) => (
                <div key={p.id} className={styles.cell}>
                  {/* eslint-disable-next-line @next/next/no-img-element */}
                  <img className={styles.cellImg} src={p.url} alt="" />
                  <button
                    className={styles.deleteBtn}
                    onClick={() => handleDeletePhoto(p.id)}
                    aria-label="Delete photo"
                  >
                    <Trash2 size={15} />
                  </button>
                </div>
              ))}
            </div>
          )}
        </Card>

        <div style={{ height: 80 }} />
      </main>

      <input ref={photoInput}  type="file" accept="image/*" hidden onChange={handlePhotoSelect} />
      <input ref={avatarInput} type="file" accept="image/*" hidden onChange={handleAvatarSelect} />

      <TabBar onFabClick={() => photoInput.current?.click()} fabLabel="Post a photo" />
    </>
  );
}

export function ProfilePage() {
  return (
    <RouteGuard>
      <ProfileScreen />
    </RouteGuard>
  );
}
